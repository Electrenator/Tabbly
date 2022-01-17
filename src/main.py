import time
import sys
import os
from browsers import Browsers
from discordPresence import DiscordPresence
from filesystem import Filesystem


class Main:
    def __init__(this, *arg, **kwargs):
        print(this, arg, kwargs)
        this.client_id = "924638024346791986"
        this.tab_logging = "log/tabs.csv"
        this.presence = DiscordPresence(this.client_id)
        this.browsers = Browsers()

    def start(this):
        try:
            while True:
                tab_count = this.browsers.countTabs()
                this.logActivity(tab_count)
                this.updateStatus(tab_count)
                time.sleep(60)
        except KeyboardInterrupt:
            pass

    def updateStatus(this, tab_count: int):
        if (tab_count > 0):
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

    def logActivity(this, tab_count: int):
        seperator = ';'

        if os.path.isfile(this.tab_logging):
            with open(this.tab_logging, "at") as logFile:
                logFile.write(
                    f"{time.time()}{seperator}{tab_count}{os.linesep}"
                )
            return
        # File does not exist yet -> make it
        Filesystem.assure_location(this.tab_logging)
        with open(this.tab_logging, "xt") as logFile:
            logFile.write(
                f"'UNIX timestamp'{seperator}'Total tab count'{os.linesep}"
            )
        # File is created -> Now really write log file
        this.logActivity(tab_count)


if __name__ == "__main__":
    app = Main(sys.argv)
    app.start()
