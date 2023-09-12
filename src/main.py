"""
Main Tabbly script. Executes the program, logs and presents it to Discord presence.
"""
import time
import sys
import os
import platform
from signal import signal, SIGTERM
from typing import Final  # Since python 3.8!

from browsers import Browsers
from discord_presence import DiscordPresence
from filesystem import assure_location, file_name_converter
from models import Setting


class Main:
    """
    Main program classified, just so less active variables are global. Also
    prevents things from executing when this file is imported as dependency
    instead of directly run.
    """

    UPDATE_CHECK_INTERVAL_SECONDS: Final[int] = 60

    def __init__(this):
        this.client_id = "924638024346791986"
        this.tab_logging = (
            "log/tabUsage"
            + (
                f"On{file_name_converter(platform.node())}"
                if len(platform.node()) > 0
                else ""
            )
            + ".csv"
        )
        this.presence = DiscordPresence(this.client_id)
        this.browsers = Browsers()
        this.previous_tab_count: None | int = None

    def start(this):
        """
        Main entry point LOOP of the program after init. Houses the looping project logic
        will run while the program is active.
        """
        try:
            while True:  # Main program loop
                this._loop()
                time.sleep(this.UPDATE_CHECK_INTERVAL_SECONDS)

        except KeyboardInterrupt:
            # Final log update before shutdown
            print("Stopping program...")
            this.log_activity(this.browsers.get_windows())
            sys.exit(0)

    def _loop(this):
        """
        This is the logic which will run with within the program loop. This function will
        repeatedly be called while the program runs.
        """
        window_data = this.browsers.get_windows()

        if window_data == this.previous_tab_count:
            return

        print("Browser change detected!")
        this.previous_tab_count = window_data

        this.update_status(sum(window_data))
        this.log_activity(window_data)

        # Flush the output to the next pipe. Would only do that on shutdown if unset
        sys.stdout.flush()

    def update_status(this, tab_count: int):
        """
        This function updates the tab use status to the given tab_count within
        presence.
        """
        if tab_count <= 0:
            print("No tabs detected")
            if this.presence.is_connected:
                this.presence.pause()
            return

        if not this.presence.is_connected:
            this.presence.resume()

        status = (
            f"Using the power of {tab_count} tab"
            + ("s" if tab_count != 1 else "")
            + " 📑"
        )
        this.presence.update(status)

    def log_activity(this, window_data: list[int]):
        """
        This function logs the browser usage activity to a csv log file specified
        within the __init__.
        """
        separator = ";"
        window_count = len(window_data)
        tab_count = sum(window_data)

        if not os.path.isfile(this.tab_logging):
            # File does not exist yet -> make it
            assure_location(this.tab_logging)
            with open(this.tab_logging, "xt", encoding="UTF-8") as log_file:
                log_file.write(
                    f"'UNIX timestamp'{separator}"
                    + f"'Total window count'{separator}"
                    + f"'Total tab count'{separator}"
                    + f"'List of total tabs per window'{os.linesep}"
                )

        with open(this.tab_logging, "at", encoding="UTF-8") as log_file:
            log_file.write(
                f"{int(time.time())}{separator}"
                + f"{window_count}{separator}"
                + f"{tab_count}{separator}"
                + f"{window_data}{os.linesep}"
            )


def exit_now(*args):
    """
    This function defines what should be done on on a Linux system shutdown. Windows already sends
    a CTRL-C event as far a I can find.
    """
    raise KeyboardInterrupt()


if __name__ == "__main__":
    signal(SIGTERM, exit_now)  # What to do on terminate request
    Setting.readFromArguments()

    app = Main()
    app.start()
