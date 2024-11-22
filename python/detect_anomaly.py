import sys
import pickle
import numpy as np
import pandas as pd

# Load the model
try:
    with open("./python/anomaly_model.pkl", "rb") as f:
        model = pickle.load(f)
except Exception as e:
    print(f"ERROR: Failed to load model: {e}")
    sys.exit(1)

# Get input metrics from command-line arguments
try:
    cpu, memory, disk_io, network_io = map(float, sys.argv[1:])
except ValueError as e:
    print(f"ERROR: Failed to convert input values: {e}")
    sys.exit(1)

# Convert inputs to DataFrame with appropriate feature names
data = pd.DataFrame([[cpu, memory, disk_io, network_io]], columns=["CPUUsage", "MemoryUsage", "DiskIO", "NetworkIO"])

# Predict anomaly using the model
try:
    prediction = model.predict(data)
    result = "Anomaly Detected" if prediction[0] == -1 else "Normal"
    print(result)
except Exception as e:
    print(f"ERROR: Failed to predict anomaly: {e}")
    sys.exit(1)
