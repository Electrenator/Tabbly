from browser.firefox import Firefox


class Browsers:
    def countTabs(this) -> int:
        return len(Firefox().getTabs() if Firefox().isRunning() else [])

    def countWindows(this) -> int:
        return len(Firefox().getWindows() if Firefox().isRunning() else [])
