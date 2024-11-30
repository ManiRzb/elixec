# Container Security Framework (elixec)

This **Container Security Framework** is a comprehensive security solution designed for containerized environments. Built with **Golang**, it integrates both static and dynamic security measures to ensure the protection of containerized applications throughout their lifecycle. This framework aims to mitigate security risks and improve the overall security posture of containers in cloud environments.

## Features

- **Image Scanning**: Scans Docker images for known vulnerabilities using tools like **Trivy**.
- **Policy Validation**: Ensures containers comply with predefined security policies (e.g., no root access, resource limits).
- **Attack Simulation**: Simulates various attack scenarios (e.g., privilege escalation, container breakouts) to identify vulnerabilities.
- **Anomaly Detection**: Uses **Isolation Forest** to detect abnormal behaviors in containers during runtime.
- **Runtime Monitoring**: Continuously monitors container metrics (CPU, memory, network activity) for real-time anomaly detection.

## Phases

The framework operates through several distinct phases:

1. **Image Scanning**: Identifies and categorizes vulnerabilities in container images.
2. **Policy Validation**: Validates container configurations against security policies, ensuring compliance.
3. **Attack Simulation**: Simulates real-world attack scenarios to evaluate container vulnerabilities.
4. **Anomaly Detection**: Monitors container behavior and uses machine learning (Isolation Forest) to detect abnormal activities.
5. **Runtime Monitoring**: Tracks container performance and flags any suspicious runtime metrics.

## How It Works

1. **Scan Image**: The framework starts by scanning the Docker image for known vulnerabilities. The severity of each vulnerability is used to assign a security score.
   
2. **Validate Policies**: It then checks the container configuration against security policies (e.g., running as root, resource limits).
   
3. **Simulate Attacks**: The framework simulates various attack scenarios (e.g., privilege escalation, container breakouts) to test container security.
   
4. **Detect Anomalies**: It monitors container metrics like CPU usage, memory, and network traffic in real-time, using anomaly detection techniques.
   
5. **Monitor Runtime**: Runtime metrics are continuously monitored to ensure the container's behavior is within expected parameters.

## Technologies Used

- **Golang**: Core implementation language.
- **Trivy**: For vulnerability scanning.
- **Isolation Forest**: Machine learning-based anomaly detection.
- **Docker**: For container deployment.

## Getting Started

### Prerequisites

- Docker
- Golang 1.16 or higher
- Trivy (for image scanning)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/ManiRzb/elixec.git
   cd elixec
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the project:
   ```bash
   go build -o elixec .
   ```

4. Run the framework:
   ```bash
   ./elixec <image-name>
   ```

### Configuration

- **Policy File**: Configure security policies by editing `configs/policies.yaml`.
- **Attack Simulations**: Customize attack scenarios in `configs/attacks.yaml`.

### Example Usage

To scan a Docker image and run the security framework:
```bash
./elixec my-image:latest
```

The framework will perform image scanning, validate policies, simulate attacks, and monitor runtime metrics, generating a security report at the end.

## Contributing

1. Fork the repository.
2. Create a new branch for your changes.
3. Commit your changes.
4. Push to the branch.
5. Create a pull request.

## License
