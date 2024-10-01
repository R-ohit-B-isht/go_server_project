import re
from datetime import datetime, timedelta

def filter_logs(log_file_path, output_file_path, hours=1):
    # Define the time range (last hour by default)
    end_time = datetime.now()
    start_time = end_time - timedelta(hours=hours)

    # Regex pattern for timestamp and "Repository not found" message
    timestamp_pattern = r'\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}'
    error_pattern = r'Repository not found'

    with open(log_file_path, 'r') as log_file, open(output_file_path, 'w') as output_file:
        for line in log_file:
            # Extract timestamp
            timestamp_match = re.search(timestamp_pattern, line)
            if timestamp_match:
                log_time = datetime.strptime(timestamp_match.group(), '%Y/%m/%d %H:%M:%S')

                # Check if the log entry is within the specified time range
                if start_time <= log_time <= end_time:
                    # Check if the line contains "Repository not found"
                    if re.search(error_pattern, line):
                        output_file.write(line)

if __name__ == "__main__":
    log_file_path = "/home/ubuntu/go_server_project/pullrequest_collection.log"
    output_file_path = "/home/ubuntu/go_server_project/filtered_logs.txt"
    filter_logs(log_file_path, output_file_path)
    print(f"Filtered logs have been written to {output_file_path}")
