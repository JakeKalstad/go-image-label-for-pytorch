# go-image-label

Reads a JSON file and organizes images into subfolders to be easily used with Pytorch `ImageFolder`

    {
        "dataFile": "birds.json",
        "imagePrefix": "images",
        "defaultPrefix": "pr_set",
        "secondaryPrefix": "pr_eggset",
        "secondaryPredicate": "egg",
        "ignorePredicate": "map",
        "outfile": "classes"
    }

dataFile is the input file - see birds.json as example format

imagePrefix is the directory that the entire image set is

defaultPrefix is where we place the subfolders and images

secondaryPrefix is where we place the subfolders and images of any data that matches our secondaryPredicate

secondaryPredicate is a string that is checked if it's contained in the image name and if so it will use the secondaryPrefix instead of the default

ignorePredicate is a string that is checked if it's contained in the image name and if so it will be discarded

outfile is the file with the unique classes listed out as a json file
