"""
Main Tabbly script. Executes the program, logs and presents it to Discord presence.
"""
import time
import sys
import os
import platform
from browsers import Browsers
from discord_presence import DiscordPresence
from filesystem import assure_location, file_name_converter


class Main:
    """
    Main program classified, just so less active variables are global. Also
    prevents things from executing when this file is imported as dependency
    instead of directly run.
    """

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

    def start(this):
        """
        Main entry point of the program after init. Houses the program loop that
        will run while the program is active.
        """
        try:
            while True:
                tab_count = this.browsers.count_tabs()
                this.update_status(tab_count)
                this.log_activity(
                    this.browsers.count_windows(),
                    tab_count,
                    this.browsers.get_windows(),
                )
                time.sleep(60)
        except KeyboardInterrupt:
            # Final log update before shutdown
            this.log_activity(
                this.browsers.count_windows(),
                this.browsers.count_tabs(),
                this.browsers.get_windows(),
            )
            sys.exit(0)

    def update_status(this, tab_count: int):
        """
        This function updates the tab use status to the given tab_count within
        presence.
        """
        if tab_count > 0:
            this.presence.resume()
            status = (
                f"Using the power of {tab_count} tab"
                + ("s" if tab_count != 1 else "")
                + " 📑"
            )
            this.presence.update(status)
        else:
            print("No tabs detected")
            this.presence.pause()

    def log_activity(this, window_count: int, tab_count: int, window_data: list[int]):
        """
        This function logs the browser usage activity to a csv log file specified
        within the __init__.
        """
        separator = ";"

        if os.path.isfile(this.tab_logging):
            with open(this.tab_logging, "at", encoding="UTF-8") as log_file:
                log_file.write(
                    f"{int(time.time())}{separator}"
                    + f"{window_count}{separator}"
                    + f"{tab_count}{separator}"
                    + f"{window_data}{os.linesep}"
                )
            return

        # File does not exist yet -> make it
        assure_location(this.tab_logging)
        with open(this.tab_logging, "xt", encoding="UTF-8") as log_file:
            log_file.write(
                f"'UNIX timestamp'{separator}"
                + f"'Total window count'{separator}"
                + f"'Total tab count'{separator}"
                + f"'List of total tabs per window'{os.linesep}"
            )

        # File is created -> Now really write log file
        this.log_activity(window_count, tab_count, window_data)


if __name__ == "__main__":
    app = Main()
    app.start()
