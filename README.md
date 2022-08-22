# wire-prod-pod

Custom Vector escape pod for regular production Vectors. No OSKR required.

## Installation

### Some general notices:

-   This will require you to put your Vector on 1.8. You cannot keep him on 1.6 if you want to use these instructions.
-   This will NOT clear user data, and the Vector SDK should work just fine.
-   It is recommended you use a Raspberry Pi 4, though you can use any Linux computer
-   For this; you do not need to pay for the voice subscription, Escape Pod, or any DDL service.

### Prerequisites

-   Raspberry Pi capable of running 64-bit Raspberry Pi OS
	-	A standard computer running Linux works too. There will be two sets of instructions for whichever you choose
-   A production Vector 1.0 that has already been setup like normal (no subscription required, he should just be active)
	-	Vector 2.0 is not supported yet
-   A computer with Google Chrome and Bluetooth
	-	If you are using Linux, go to [chrome://flags](chrome://flags) and enable `Experimental web platform features`, then relaunch Chrome

### Set up the bot

1.	Put your Vector on the charger and hold the button for 15 seconds. The light will turn off after ~5 seconds, that is normal. Just keep holding until he turns back on (do *NOT* clear user data!).

2.	He should be at a screen that shows `anki.com/v`. On a computer with Bluetooth support (preferably Windows or macOS), go to [https://www.project-victor.org/noflow-devsetup](https://www.project-victor.org/noflow-devsetup) in a Chrome browser and connect to Vector like it tells you to. Do not check `Enable auto-setup flow`.
	-	You may need to reload the page a few times for it to connect correctly
	-	If you are using Linux, you may need to open up Bluetooth settings and keep it discovering in the background

3.	You should now be at a terminal-like interface. In that interface, connect Vector to Wi-Fi with this command: 

```
wifi-connect "ssid" "password"
```
-	Replace `ssid` with your network name, `password` with the network password
-	Example: `wifi-connect "AnkiRobits" "KlaatuBaradaNikto!"`

4.	Enter the following command into the site: 
```
ota-start http://173.20.162.183:81/escapepod-prod-1.8.ota
```

5.	Wait for that to finish. Once he has rebooted, continue on to the next set of instructions.

### Set up wire-pod

Only do one of the following sets of instructions.

#### Option one: Raspberry Pi 4

1. Download and install [Raspberry Pi Imager](https://www.raspberrypi.com/software/), insert a micro-SD card into your computer/into a reader

2. Follow along with [this video](https://www.youtube.com/watch?v=PUNcRqgt1cA). This is an example of how I would set it up for my network/liking. The only things you should change are the user password (located under `pi`, change it to anything you want and make sure you remember it) and the network information (you should put your network name instead of `AnkiRobits`, your network password instead of `KlaatuBaradaNikto!`)

3. Insert the SD card into your Pi, wait a few minutes for it to boot up

4. On your computer, open up a Terminal (or Powershell on Windows) and run the following command:

```
ssh pi@escapepod
```

5. If it gives you a prompt, enter `yes`. Then enter in the user password you chose in Raspberry Pi Imager (under `pi`)

6. The terminal should show `pi@escapepod ~ $`. If you are there, run the following command:

```
wget -O - https://wire.my.to/setup-wire-pod.sh | bash
```

7. After that completes, open a browser and go to [http://escapepod:8080](http://escapepod:8080). From there, click on "Set up wire-pod" and do as it says. You do not need to do the part at the bottom where it says `Choose file`, that is reserved for OSKR/dev bots.

8. Voice commands should now work!

#### Option two: (almost) Any Linux system

(Your distribution must have either pacman, dnf, or apt, make sure you have ports 443 and 8080 open)

1. Open a terminal

2. Use the following commands to change your hostname to `escapepod` and have that register with the network (this will change your computer's name):

```
sudo hostnamectl set-hostname escapepod
systemctl restart avahi-daemon
systemctl enable avahi-daemon
```

(ignore any errors that may pop up, they are normal)

3. Open a terminal and run this command:

```
wget -O - https://wire.my.to/setup-wire-pod.sh | bash
```

4. Once that has completed, go to the URL it tells you to go to in a browser and click on `Set up wire-pod (API keys, STT service, etc)`. From there, follow the instructions. It should then be set up.

## Updating

Wire-pod auto-updates once a day. To force an update, run the following commands:

```
cd ~/wire-prod-pod
sudo ./update.sh
```

## Web interface

Chipper hosts a web interface at port 8080. This can be used to create custom intents and to configure specific bots.

To get to it, open a browser and go to `http://serverip:8080`, replacing serverip with the IP address of the machine running the chipper server. If you are running the browser on the machine running chipper, you can go to `http://localhost:8080`

- Set up wire-pod
	- This is where you would enter your API keys for [Houndify](https://www.houndify.com/dashboard) and [WeatherAPI](https://www.weatherapi.com/) if you want to set those up.
- Configure user settings
	- Custom intents
		- Example: You want to create a custom intent that allows Vector to turn the lights off. The transcribed text that matches to this intent should include "lights off" and other variations like "lid off" for better detection. It will execute a python script located in your user directory called `vlight.py`. It should be launched with the `off` variable because the lights are being turned off. This script turns the lights off and connects to Vector so he says "The lights are off!". You have multiple bots registered with the SDK so a serial number must be specified. After the SDK program is complete, chipper should send `intent_greeting_goodnight`. The following screenshot is a correct configuration for this case. The `Add intent` button would be pressed after everything is put in.
		- ![Custom Intent Screenshot](./images/customIntent.png)
		- (If `!botSerial` is put into the program arguments, chipper will substitute it for the serial number of the bot that is making a request to it.)
	- Bot configurations
		- Example: You have a bot on 1.6 with a serial number of 0060059b; located in Des Moines, Iowa; and you would like to use Fahrenheit for the temperature unit in the weather command. The following screenshot is a correct configuration for that bot. Once that is input, you would press "Add bot". It would then be registered with chipper.
		- ![Bot Configuration Screenshot](./images/botConfig.png)
- See logs
	- This shows the logs produced by wire-pod. Updates every 200ms, shows the last 30 lines.

## Status

OS Support:

- Arch
- Debian/Ubuntu/other APT distros
- Fedora/openSUSE
- Windows under WSL

Architecture support:

- amd64/x86_64
- arm64/aarch64
- armv7l/arm32 (picovoice only)

Things wire-prod-pod has worked on:

- Raspberry Pi 4B+ 4GB RAM with Raspberry Pi OS
	- Must be 64-bit OS if using Coqui
- Raspberry Pi 4B+ 4GB RAM with Manjaro 22.04
- Nintendo Switch with L4T Ubuntu
- Desktop with Ryzen 5 3600, 16 GB RAM with Ubuntu 22.04
- Laptop with mobile i7
- Desktop with Xeon X5650

General notes:

- If you get this error when running chipper, you are using a port that is being taken up by a program already: `panic: runtime error: invalid memory address or nil pointer dereference`
- You have to speak a little slower than normal for Coqui STT to understand you.

Current implemented actions:

- Good robot
- Bad robot
- Change your eye color
- Change your eye color to <color>
	- blue, purple, teal, green, yellow
- How old are you
- Start exploring ("deploring" works better)
- Go home (or "go to your charger")
- Go to sleep
- Good morning
- Good night
- What time is it
- Goodbye
- Happy new year
- Happy holidays
- Hello
- Sign in alexa
- Sign out alexa
- I love you
- Move forward
- Turn left
- Turn right
- Roll your cube
- Pop a wheelie
- Fistbump
- Blackjack (say yes/no instead of hit/stand)
- Yes (affirmative)
- No (negative)
- What's my name
- Take a photo
- Take a photo of me
- What's the weather
	- Requires API setup
- What's the weather in <location>
	- Requires API setup
- Im sorry
- Back up
- Come here
- Volume down
- Be quiet
- Volume up
- Look at me
- Set the volume to <volume>
	- High, medium high, medium, medium low, low
- Shut up
- My name is <name>
- I have a question
	- Requires API setup
- Set a timer for <time> seconds
- Set a timer for <time> minutes
- Check the timer
- Stop the timer
- Dance
- Pick up the cube
- Fetch the cube
- Find the cube
- Do a trick
- Record a message for <name>
	- Enable `Messaging` feature in webViz Features tab
- Play a message for <name>
	- Enable `Messaging` feature in webViz Features tab

## Credits

- [Digital Dream Labs](https://github.com/digital-dream-labs) for saving Vector and for open sourcing chipper which made this possible
- [dietb](https://github.com/dietb) for rewriting chipper and giving tips
- [GitHub Copilot](https://copilot.github.com/) for being awesome
- [xanathon](https://github.com/xanathon) for making the web interface look nice and for the publicity
