import time
import sys
from browserTab import BrowserTabs

from discordPresence import DiscordPresence


class Main:
    def __init__(self, *arg, **kwargs):
        print(self, arg, kwargs)
        self.client_id = "924638024346791986"
        self.presence = DiscordPresence(self.client_id)
        self.tabs = BrowserTabs()

    def start(self):
        try:
            while True:
                status = f"Using the power of {self.tabs.count()} tabs 📑"
                print(self.presence.update(status))
                time.sleep(1)
        except KeyboardInterrupt:
            pass


if __name__ == "__main__":
    app = Main(sys.argv)
    app.start()
