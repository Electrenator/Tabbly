"""
Utility functions needed for whatever reason go into this file.
"""
from argparse import ArgumentParser, Namespace


def get_environment_arguments() -> Namespace:
    """
    Creates this projects argument parser.

    Returns: The parser used to get the projects arguments
    """
    parser = ArgumentParser(
        description="A program for showing your browser tab usage within Discord's rich presence "
        + " and to log that usage.",
        epilog="Project source can be found over at; https://github.com/Electrenator/Tabbly",
        allow_abbrev=False,
    )
    parser.add_argument(
        "-v", "--verbose", action="store_true", help="show extended console output."
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Does't actually write outputs to a file or Discord.",
    )
    return parser.parse_args()
