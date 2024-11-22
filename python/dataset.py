import pandas as pd
import numpy as np
import random

# Configuration for the dataset
NUM_SAMPLES = 50000
ANOMALY_PROB = 0.1  # Probability of an anomaly in the dataset

def generate_metrics(is_anomalous=False):
    """
    Generate a single row of container metrics.
    """
    if is_anomalous:
        cpu = np.random.uniform(80, 100)  # High CPU usage
        memory = np.random.uniform(80, 100)  # High memory usage
        disk_io = np.random.uniform(500, 1000)  # High disk I/O
        network_io = np.random.uniform(500, 1000)  # High network I/O
        label = 1  # Anomalous
    else:
        cpu = np.random.uniform(0, 70)  # Normal CPU usage
        memory = np.random.uniform(0, 70)  # Normal memory usage
        disk_io = np.random.uniform(50, 200)  # Normal disk I/O
        network_io = np.random.uniform(50, 200)  # Normal network I/O
        label = 0  # Normal

    return [cpu, memory, disk_io, network_io, label]

# Generate the dataset
data = []
for _ in range(NUM_SAMPLES):
    is_anomalous = random.random() < ANOMALY_PROB
    data.append(generate_metrics(is_anomalous))

# Create a DataFrame
df = pd.DataFrame(data, columns=["CPUUsage", "MemoryUsage", "DiskIO", "NetworkIO", "Label"])

# Save the dataset to a CSV file
df.to_csv("./python/container_metrics.csv", index=False)
print("Dataset generated and saved as container_metrics.csv")
