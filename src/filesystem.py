"""
This module stores all functions related to filesystem operations written for Tabbly.
"""
import glob
import os
from re import sub


def find_files(path: str) -> list:
    """
    Finds files within given path while allowing relative linux like glob paths to be used.

    Replaces `~` with the users home directory, rest of the glob syntax is explained within the
    documentation; https://docs.python.org/3.8/library/glob.html
    """
    if path.startswith("~"):
        path = path.replace("~", os.path.expanduser("~"), 1)

    return glob.glob(path)


def assure_location(file_path: str, dry_run: bool = False):
    """
    Assures a directory location for the given file path without creating the file
    itself. Does nothing when the file path already exists but can run without actually
    creating any directoires with `dry_run`.

    With the file_path `../nonExistantDir/subdir/newFile.extention` the directories
    `nonExistantDir/subdir` within the current path's parrent direcory should be created
    but the `newFile.extention` should not be made.
    """
    file_path = os.path.realpath(file_path)
    file_name = os.path.basename(file_path)
    parent_dir_path = file_path.replace(file_name, "")

    if not dry_run:
        try:
            os.makedirs(parent_dir_path)
        except FileExistsError:
            pass
    else:
        print(f"Would have created directory '{parent_dir_path}'")


def file_name_converter(name_string: str):
    """
    Returns a UpperCamelCase version of a given name_string to use as for file names.
    """
    return sub(r"(-|_)+", "", name_string.title())
