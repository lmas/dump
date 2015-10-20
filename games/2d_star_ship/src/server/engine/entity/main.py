

class EntityManager(object):
    # Heavily based on code from:
    # https://github.com/seanfisk/ecs/blob/master/ecs/managers.py
    def __init__(self):
        self.eid = 0
        self.entities = {}
        self.systems = {}


    def add_system(self, func, comp_type):
        if comp_type not in self.systems:
            self.systems[comp_type] = []
        self.systems[comp_type].append(func)

    def del_system(self, func):
        for comp_type in list(self.systems.keys()):
            try:
                self.systems[comp_type].remove(func)
                if self.systems[comp_type] == []:
                    del self.systems[comp_type]
            except KeyError:
                pass

    def update_systems(self, delta_time=0):
        for comp_type, systems in self.systems.items():
            entities = self.entities_for_component(comp_type)
            for func in systems:
                func(self, entities, delta_time)


    def add_entity(self):
        eid = self.eid
        self.eid += 1
        return eid

    def del_entity(self, eid):
        for comp_type in list(self.entities.keys()):
            try:
                del self.entities[comp_type][eid]
                if self.entities[comp_type] == {}:
                    del self.entities[comp_type]
            except KeyError:
                pass

    def add_component(self, eid, component):
        comp_type = type(component)
        if comp_type not in self.entities:
            self.entities[comp_type] = {}
        self.entities[comp_type][eid] = component

    def del_component(self, eid, comp_type):
        try:
            del self.entities[comp_type][eid]
            if self.entities[comp_type] == {}:
                del self.entities[comp_type]
        except KeyError:
            pass

    def get_component(self, eid, comp_type):
        try:
            return self.entities[comp_type][eid]
        except KeyError:
            return None

    def entities_for_component(self, comp_type):
        try:
            return self.entities[comp_type].keys()
        except KeyError:
            return {}

    def components_for_entity(self, eid):
        tmp = []
        for entities in self.entities.values():
            try:
                tmp.append(entities[eid])
            except KeyError:
                pass
        return tmp

