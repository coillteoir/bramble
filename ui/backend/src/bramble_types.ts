export class TaskSpec {
  image: string;
  command: string[];
  dependencies?: string[];

  constructor(image: string, command: string[], dependencies?: string[]) {
    this.image = image;
    this.command = command;
    this.dependencies = dependencies;
  }
}

export class PLtask {
  name: string;
  spec: TaskSpec;

  constructor(name: string, spec: TaskSpec) {
    this.name = name;
    this.spec = spec;
  }
}

export class TaskRef {
  name: string;
  dependencies?: string[];
  constructor(name: string, dependencies?: string[]) {
    this.name = name;
    this.dependencies = dependencies;
  }
}

export class Pipeline {
  metadata: {
    name: string;
    namespace: string;
  };
  spec: {
    tasks?: PLtask[];
    taskRefs?: TaskRef[];
  };

  constructor(
    metadata: {
      name: string;
      namespace: string;
    },
    spec: {
      tasks?: PLtask[];
      taskRefs?: TaskRef[];
    },
  ) {
    this.metadata = metadata;
    this.spec = spec;
  }
}
