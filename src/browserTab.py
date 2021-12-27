
class BrowserTabs:
    def __init__(self):
        self.fakeTabs = 0

    def count(self) -> int:
        self.fakeTabs += 1
        return self.fakeTabs

