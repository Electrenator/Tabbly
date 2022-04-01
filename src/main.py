import time
import sys
import os
from browsers import Browsers
from discord_presence import DiscordPresence
import filesystem


class Main:
    def __init__(this, *arg, **kwargs):
        print(this, arg, kwargs)
        this.client_id = "924638024346791986"
        this.tab_logging = "log/tabUsage.csv"
        this.presence = DiscordPresence(this.client_id)
        this.browsers = Browsers()

    def start(this):
        try:
            while True:
                tab_count = this.browsers.count_tabs()
                this.update_status(tab_count)
                this.log_activity(this.browsers.count_windows(), tab_count)
                time.sleep(60)
        except KeyboardInterrupt:
            # Final log update before shutdown
            this.log_activity(this.browsers.count_windows(),
                              this.browsers.count_tabs())
            exit(0)

    def update_status(this, tab_count: int):
        if tab_count > 0:
            this.presence.resume()
            status = (
                f"Using the power of {tab_count} tab"
                + ("s" if tab_count != 1 else "")
                + " 📑"
            )
            print(this.presence.update(status))
        else:
            print("No tabs detected")
            this.presence.pause()

    def log_activity(this, window_count: int, tab_count: int):
        seperator = ';'

        if os.path.isfile(this.tab_logging):
            with open(this.tab_logging, "at", encoding="UTF-8") as log_file:
                log_file.write(
                    f"{int(time.time())}{seperator}" +
                    f"{window_count}{seperator}" +
                    f"{tab_count}{os.linesep}"
                )
            return
        # File does not exist yet -> make it
        filesystem.assure_location(this.tab_logging)
        with open(this.tab_logging, "xt", encoding="UTF-8") as log_file:
            log_file.write(
                f"'UNIX timestamp'{seperator}" +
                f"'Total window count'{seperator}" +
                f"'Total tab count'{os.linesep}"
            )
        # File is created -> Now really write log file
        this.log_activity(window_count, tab_count)


if __name__ == "__main__":
    app = Main(sys.argv)
    app.start()
