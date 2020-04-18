import sys
import os
import yaml, json
from typing import List


def isYaml(file: str) -> bool:
    # Check if the file has a yaml extension
    #
    # Parameters:
    #   file (str): target filename
    #
    # Returns:
    #   isYaml (bool)

    _, fileExt = os.path.splitext(file)
    return fileExt.lower() == ".yaml" or fileExt.lower() == ".yml"

def listdir_fullpath(d: str) -> List(str):
    # A better listdir function which returns the full path for list dir function
    #
    # Parameters:
    #   d (str): target dirname
    #
    # Returns:
    #   listOfFileNames (List(str)): List of file names with full path

    return [os.path.join(d, f) for f in os.listdir(d)]

if __name__ == "__main__":
    '''A cli to convert all the YAML files on dir to JSON files'''
    if len(sys.argv) != 2:
        print("Please provide the dir containing the yaml test files")
        sys.exit(1)
    yamlFiles = list(filter(isYaml, listdir_fullpath(sys.argv[1])))
    for yamlFile in yamlFiles:
        with open(yamlFile) as file:
            target = yaml.load(file)
        filename, _ = os.path.splitext(yamlFile)
        with open(filename + ".json", "w+") as jsonFile:
            json.dump(target, jsonFile)
