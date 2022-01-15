import time
import sys
from browsers import Browsers

from discordPresence import DiscordPresence


class Main:
    def __init__(this, *arg, **kwargs):
        print(this, arg, kwargs)
        this.client_id = "924638024346791986"
        this.presence = DiscordPresence(this.client_id)
        this.browsers = Browsers()

    def start(this):
        try:
            while True:
                tab_count = this.browsers.countTabs()
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


if __name__ == "__main__":
    app = Main(sys.argv)
    app.start()
