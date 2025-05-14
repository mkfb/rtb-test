import asyncio
import aiohttp
import time
import random
import sys

# Define the target URL.  This should be the /bid endpoint of your Nginx proxy.
TARGET_URL = "http://nginx/bid"  # Replace with your actual URL

# Define the number of requests per second (RPS).
RPS = 100
# Number of concurrent requests.  Adjust this based on your system's capacity.
CONCURRENCY = 20
#  Timeout for each individual request.
TIMEOUT = 5  # Seconds

# Define a sample OpenRTB 2.5 bid request.  This is a simplified
# version; a real request would be much more complex.  It's crucial
# to have a valid, albeit dummy, OpenRTB request here.
SAMPLE_REQUEST = {
    "id": "request-" + str(random.randint(100000, 999999)),
    "imp": [
        {
            "id": "imp-123",
            "banner": {
                "w": 300,
                "h": 250,
            },
            "displaymanager": "openrtb-sim",
            "displaymanagerver": "1.0",
            "tagid": "tag-456",
            "bidfloor": 0.5,
            "bidfloorcur": "USD",
        }
    ],
    "device": {
        "ua": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
        "ip": "192.168.1.100",
        "os": "Mac OS X",
        "model": "MacBook Pro",
        "connectiontype": 2,
    },
    "user": {
        "id": "user-789",
    },
    "tmax": 1000, # Set a timeout
}


# Function to send a single request.  This function is an async
# function, meaning it can be run concurrently with other requests.
async def send_request(session: aiohttp.ClientSession):
    """
    Sends a single OpenRTB bid request to the target URL using aiohttp.
    Handles potential errors and logs the response status.

    Args:
        session: The aiohttp ClientSession to use for the request.  This
                 is passed in to allow for connection pooling.
    """
    try:
        # Use a timeout for the request to prevent it from hanging indefinitely.
        async with session.post(
            TARGET_URL,
            json=SAMPLE_REQUEST,
            timeout=TIMEOUT,
        ) as response:
            # Log the response status code.  This is essential for monitoring
            # the server's performance and identifying any errors.
            # print(f"Request to {TARGET_URL} completed with status: {response.status}")
            # Read the response body (optional).  This is useful for debugging
            # or if you need to validate the response content.  However, in a
            # high-throughput scenario, you might want to skip reading the body
            # to reduce overhead.
            response_text = await response.text()
            print(f"Response: {response_text}")

            # Raise an exception for bad status codes (400, 500, etc.).
            response.raise_for_status()  # Raise exception for bad status
    except aiohttp.ClientError as e:
        # Handle client-side errors (e.g., connection errors, invalid URL).
        print(f"Error sending request to {TARGET_URL}: {e}")
    except asyncio.TimeoutError:
        # Handle timeouts separately.
        print(f"Request to {TARGET_URL} timed out after {TIMEOUT} seconds")
    except Exception as e:
        # Handle any other unexpected errors.  This is a catch-all
        # to prevent the program from crashing.
        print(f"An unexpected error occurred: {e}")



# Function to generate traffic at the specified RPS.
async def generate_traffic():
    """
    Generates traffic to the target URL at the specified RPS using asyncio.
    It calculates the required delay between requests to achieve the target RPS
    and uses aiohttp for asynchronous HTTP requests.
    """
    # Calculate the delay between requests to achieve the target RPS.
    # The delay is in seconds.
    delay = 1 / RPS
    print(f"Sending requests to {TARGET_URL} at {RPS} RPS (delay: {delay:.6f} seconds)")

    # Create an aiohttp ClientSession.  This session will be used for all
    # requests, allowing for connection pooling and improved performance.
    async with aiohttp.ClientSession() as session:
        while True:
            # Create a list of tasks to send concurrently.  We limit the number
            # of concurrent tasks to CONCURRENCY to avoid overwhelming the server
            # or our own system.
            tasks = [send_request(session) for _ in range(CONCURRENCY)]
            # Use asyncio.gather to run the tasks concurrently.  This will send
            # the requests as fast as possible, up to the concurrency limit.
            await asyncio.gather(*tasks)
            # Introduce a delay to achieve the target RPS.  We use asyncio.sleep
            # to avoid blocking the event loop.
            await asyncio.sleep(delay * CONCURRENCY)



# Main entry point of the script.
if __name__ == "__main__":
    # Check if the target URL is provided as a command-line argument.
    # If it is, update the TARGET_URL.
    if len(sys.argv) > 1:
        TARGET_URL = sys.argv[1]
        print(f"Using target URL from command line: {TARGET_URL}")

    # Run the traffic generation function using asyncio.run.  This sets up
    # the asyncio event loop and runs the generate_traffic function until it completes.
    asyncio.run(generate_traffic())
