import sys
import os
import yaml, json

def isYaml(file: str):
    _, fileExt = os.path.splitext(file)
    return fileExt.lower() == ".yaml" or fileExt.lower() == ".yml"

def listdir_fullpath(d):
    return [os.path.join(d, f) for f in os.listdir(d)]

if __name__ == "__main__":
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
