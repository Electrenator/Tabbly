from browser.firefox import Firefox


class Browsers:
    def countTabs(this) -> int:
        return len(Firefox().getTabs() if Firefox().isRunning() else [])
