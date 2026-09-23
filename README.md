# SSGO - High School Subject Selection Online System

SSGO (Subject Selection Go) is a robust and responsive system designed to facilitate high school elective subject selections for students, teachers, and administrators. It comprises a Go-based backend server and a modern Vue 3 frontend application.

---

## 🛠️ Environment Configuration (`.env`)

Before running or deploying the application, you must configure your environment variables. 

### 1. Setup Environment File
Copy the provided `.env.example` template into a new `.env` file in the project root:
```bash
cp .env.example .env
```

### 2. Required Configurations
Open `.env` and fill in the following secure credentials (which are left empty in the template):
* `SUPER_ADMIN_PASSWORD`: Set a secure password for the Configuration Dashboard.
* `JWT_PRIVATE_KEY`: Set a complex, secret string used for signing authentication JSON Web Tokens (JWTs).

### 3. Environment Variables Directory

| Variable Name | Description | Default / Example Value |
| :--- | :--- | :--- |
| `SUPER_ADMIN_USERNAME` | Username for the Configuration Dashboard. | `root` |
| `SUPER_ADMIN_PASSWORD` | Password for the Configuration Dashboard. | *[Required; User-defined]* |
| `JWT_PRIVATE_KEY` | Secret key used to sign JWT authentication tokens. | *[Required; User-defined]* |
| `DEFAULT_SCHOOL_NAME` | The school name in Chinese. | `聖公會李炳中學` |
| `DEFAULT_SCHOOL_WEBSITE`| URL to the school's website. | `https://liping.edu.hk` |
| `DEFAULT_SYSTEM_TITLE` | The title displayed in the header for formal select. | `高中選科系統` |
| `DEFAULT_MOCK_SYSTEM_TITLE`| The title displayed in the header for mock select. | `高中模擬選科系統` |
| `DEFAULT_INTRODUCTION_MAKRDOWN`| Markdown text shown under "注意" (Attention) during formal select. | *(Detailed Chinese markdown guidelines)* |
| `DEFAULT_MOCK_INTRODUCTION_MAKRDOWN`| Markdown text shown under "注意" during mock select. | *(Mock-specific Chinese guidelines)* |
| `DEFAULT_INSTRUCTION_MARKDOWN`| Markdown text shown under "使用須知" (Instructions) for selections. | `1. 將在「尚未編排的選科組合」...` |
| `DEFAULT_NOT_ACCEPT_MARKDOWN` | Text shown to students when formal selection is closed. | `選科程序已經截止。如有任何查詢...` |
| `DEFAULT_MOCK_NOT_ACCEPT_MARKDOWN`| Text shown to students when mock selection is closed. | `模擬選科程序已經截止...` |
| `DEFAULT_CONFIRM_MARKDOWN` | Text shown in the parent signature confirmation modal. | `是否確定簽署及選科次序？一經確定後...` |
| `DEFAULT_SUBJECTS_JSON` | JSON array representing all elective subjects. | `[{"code":"bio", "group":1, ...}]` |
| `DEFAULT_COMBINATIONS_JSON` | JSON array representing all valid subject combinations. | `[{"id":0, "subjects":["bio","chem"]}]` |

---

## 🚀 Deployment Guide

SSGO is packed with convenient shell scripts to manage compilation and containerized deployment.

### Production Deployment (using Docker)

1. **Configure Environment Variables**:
   Follow the **Environment Configuration** steps above to create and configure your `.env` file.

2. **Build the Application**:
   Run the build script in the root directory. This script compiles the Vue 3 frontend assets and copies them into the Go server, then builds the Go binary executable:
   ```bash
   ./build.sh
   ```

3. **Start the System**:
   Start the application in containerized mode (via Docker Compose):
   ```bash
   ./start.sh
   ```
   *Note: This script launches the containers in detached (background) mode.*

4. **Monitoring and Management Utilities**:
   - **Stop Services**: Stop and clean up active containers by running `./stop.sh`.
   - **Restart Services**: Automatically stop and start the services again by running `./restart.sh`.

---

## 💻 Local Development (No Docker)

For developers wishing to run and debug the backend and frontend locally without Docker:

1. Make sure `.env` is fully set up in the root directory.
2. Compile and run both the Go backend and launch local hot-reloading for the frontend:
   ```bash
   ./run_local.sh
   ```
   *Note: Ensure you have `Go` and `NodeJS` installed on your local development machine.*