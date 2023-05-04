# Tabbly
A program for showing your tab usage within Discord's rich presence and to log that usage. Currently only supports usage on Firefox within GNU/Linux.

![Tabbly in use within Discord. The application is displaying the usage of 83 active browser tabs.](https://user-images.githubusercontent.com/18311389/151074155-78ccf239-5127-4e7a-8380-f7038ade6338.png)

## How to run
You will firstly need to have the following things installed and have access to a terminal within the project's folder.
- Python 3.8+
- Python 3.8+ PIP

### TL;DR
Use the following command **And** make sure `Display current activity as a status message.` is enabled within the privacy settings!
```bash
python -m venv venv && source venv/bin/activate && pip install -r requirements.txt && python src/main.py; deactivate
```

### Full explanation
After installing that, you will need a virtual environment for python to run in. This can be created and entered with the following command on GNU/Linux. It may be necessary to type `python3` instead of `python` if you also have Python version 2 installed.
```bash
python -m venv venv
source venv/bin/activate
```
Note: this virtual environment can be deactivated with the `deactivate` command and cleaned up after by just removing the `venv` folder.

Finally you will need to install the dependencies of this program. This can be done with the following command but `pip` may need to be replaced with `pip3` depending on what's installed on your PC.
```bash
pip install -r requirements.txt
```
Now you can run `python src/main.py` to start running the program. The only thing you need to make sure of now is that the `Display current activity as a status message.` setting is enabled within Discord. Tabbly won't show anything if it isn't.


Note: When you want to contribute, you should probably also add the dev dependencies so pylint and black can be used during development. This can be done in the venv described above using the following command.
```bash
pip install -r requirements.dev.txt
```
