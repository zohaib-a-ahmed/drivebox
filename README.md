# DriveBox

DriveBox is a command-line interface (CLI) tool designed to simplify managing files on Google Drive. It allows users to upload, download, and organize their files directly from the terminal, providing a streamlined experience for interacting with Google Drive without the need for a graphical interface. The goal for this package is to manipulate and mesh local and Drive file systems seamlessly via the Drive API.

## Features

- **Upload Files**: Easily upload files to your Google Drive.
- **Download Files**: Download files from your Google Drive to your local system.
- **Manage Directories**: Create and search for directories within your Google Drive.

## Setup Instructions

Before using DriveBox, you'll need to set up a Google Cloud project and obtain your OAuth 2.0 credentials (Client ID and Client Secret). Follow these steps to get started:

1. **Create a Google Cloud Project**:

   - Visit the [Google Cloud Console](https://console.cloud.google.com/).
   - Create a new project or select an existing one.

2. **Enable the Drive API**:

   - Navigate to the "APIs & Services > Dashboard" section.
   - Click "+ ENABLE APIS AND SERVICES" and search for "Google Drive API". Enable it.

3. **Create Credentials**:

   - Go to "APIs & Services > Credentials" and click "Create Credentials".
   - Choose "OAuth client ID", then select "Desktop app" as the application type.
   - Save the provided Client ID and Client Secret.

4. **Configure OAuth Consent Screen**:

   - Set up the OAuth consent screen in the "Credentials" section by providing the necessary information about your application.

5. **Create a `.env` File**:
   - In the root directory of DriveBox, create a `.env` file containing your Client ID and Client Secret:
     ```
     GOOGLE_CLIENT_ID=your_client_id_here
     GOOGLE_CLIENT_SECRET=your_client_secret_here
     ```

## Usage

### Authentication

Before using DriveBox to manage your files, you must authenticate with Google Drive:

```sh
drivebox auth in
```

Follow the prompts to sign in with your Google account.

### Uploading Files

To upload a file to Google Drive:

```sh
drivebox upload <path_to_file>
```

Optionally, specify/create a parent directory in Google Drive to upload to:

```sh
drivebox upload parent <path_to_file>
```

### Downloading Files

To download a file from Google Drive:

```sh
drivebox unload <file_name> <optional_path_destination>
```

## Development

- Clone the repository: `git clone https://github.com/zohaib-a-ahmed/drivebox.git`
- Navigate to the cloned directory: `cd drivebox`

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues to suggest improvements or add new features.
