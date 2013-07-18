#!/usr/bin/env python

import pkgutil
import sys

# Source: http://stackoverflow.com/questions/1057431/loading-all-modules-in-a-folder-in-python/8556471#8556471
# I made some slight changes to the code..

def import_all_modules(dirname):
    modules = []
    for importer, package_name, _ in pkgutil.iter_modules([dirname]):
        full_package_name = '%s.%s' % (dirname, package_name)
        if full_package_name not in sys.modules:
            module = importer.find_module(package_name
                        ).load_module(full_package_name)
            modules.append(module)
    return modules
