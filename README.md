# lakitu

![lakitu version](https://img.shields.io/badge/version-0.1.0-lightgreen.svg)
![GPL v3 License](https://img.shields.io/badge/license-GPL%20v3-blue.svg)

![lakitu UI](https://user-images.githubusercontent.com/14242625/141721197-b45a707a-f5e8-4a13-9ad9-9bde929cf40f.gif)

Manage your cloud gaming EC2 instance wiith lakitu. 

<p align="center">
	<img src="https://user-images.githubusercontent.com/14242625/141721197-b45a707a-f5e8-4a13-9ad9-9bde929cf40f.gif">
</p>

## Features

- Create and delete cloud gaming machine on AWS
- Start/stop (save) your cloud machine using images and snapshots
- Automatically set a password for login to cloud machine
- Transfer your cloud machine between regions

## Usage

Download the [server](https://github.com/sereneblue/lakitu/releases) from the releases page, run, and visit ``localhost:8080``

Follow the initial setup guide to add your AWS secret key and access key so that lakitu can work.

After creating a machine, download the lakitu-cli binary to your cloud instance and run the bootstrap command like below to start setting up your cloud machine using the [Parsec Cloud Preperation Tool](https://github.com/parsec-cloud/Parsec-Cloud-Preparation-Tool).

**NOTE: Windows Defender may detect the lakitu-cli binary as malware. This is a false positive.**

```
	# Will install Parsec using the Parsec Cloud Preparation Tool
	> lakitu-cli bootstrap
```

You can also use lakitu-cli to initialize instance stores and volumes or create and attach a volume from a snapshot ID.

```
	# this will initialize and format an offline instance store
	> lakitu-cli mount

	# this will initialize and format an attached volume and do the above 
	> lakitu-cli mount new

	# this will create and attach a volume from the provided snapshot ID and mount instance stores
	# NOTE: An IAM role needs to be attached to the instance for this to work, this is automatically done using the server
	> lakitu-cli mount snap-1234567890
```

lakitu-cli can be used as a standalone utility without the server. If you do plan to use lakitu, you'll need to add lakitu-cli to your PATH to automatically mount storage when a machine is started.

### Note for developers

lakitu is built with SvelteKit and Echo. To start the dev server, run the below commands:

```
	$ cd web
	$ npm install
	$ npm run dev
	$ cd ..
	$ go run cmd/server/main.go
```

## Bugs

Please feel free to submit a pull request or fork this project. Unfortunately, I do not plan to work on this any further so do keep that in mind if you encounter issues.

## Disclaimer

This project is a prototype and not affiliated with Amazon or Parsec. lakitu manages snapshots, images and instances on AWS and there are probably bugs that can cause some issues. Use at your own risk.