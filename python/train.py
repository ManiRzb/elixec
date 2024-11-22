import pandas as pd
from sklearn.ensemble import IsolationForest
from sklearn.model_selection import train_test_split
from sklearn.metrics import classification_report
import pickle

# Load the dataset
data = pd.read_csv("./python/container_metrics.csv")

# Features and labels
X = data[["CPUUsage", "MemoryUsage", "DiskIO", "NetworkIO"]]
y = data["Label"]

# Split into training and test sets
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

# Train an Isolation Forest model
model = IsolationForest(contamination=0.1, random_state=42)
model.fit(X_train)

# Predict anomalies
y_pred = model.predict(X_test)
y_pred = [1 if x == -1 else 0 for x in y_pred]  # Convert -1 (anomaly) to 1, 1 (normal) to 0

# Print evaluation metrics
print("Classification Report:")
print(classification_report(y_test, y_pred))

# Save the trained model
with open("./python/anomaly_model.pkl", "wb") as f:
    pickle.dump(model, f)
print("Model trained and saved as anomaly_model.pkl")
