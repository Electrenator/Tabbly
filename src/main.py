import time
import sys
from browser import Browser

from discordPresence import DiscordPresence


class Main:
    def __init__(this, *arg, **kwargs):
        print(this, arg, kwargs)
        this.client_id = "924638024346791986"
        this.presence = DiscordPresence(this.client_id)
        this.browser = Browser()

    def start(this):
        try:
            while True:
                status = f"Using the power of {this.browser.countTabs()} tabs 📑"
                print(this.presence.update(status))
                time.sleep(15)
        except KeyboardInterrupt:
            pass


if __name__ == "__main__":
    app = Main(sys.argv)
    app.start()
