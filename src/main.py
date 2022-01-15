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
                status = (
                    f"Using the power of {tab_count} tab"
                    + ("s" if tab_count != 1 else "")
                    + " 📑"
                )
                print(this.presence.update(status))
                time.sleep(60)
        except KeyboardInterrupt:
            pass


if __name__ == "__main__":
    app = Main(sys.argv)
    app.start()
