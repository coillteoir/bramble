export class TaskSpec {
  image: string;
  command: string[];

  constructor(image: string, command: string[]) {
    this.image = image;
    this.command = command;
    console.log("TaskSpec: ", this);
  }
}

export class PLtask {
  name: string;
  spec: TaskSpec;
  dependencies?: string[];

  constructor(name: string, spec: TaskSpec, dependencies?: string[]) {
    this.name = name;
    this.spec = spec;
    this.dependencies = dependencies;
    console.log("PLtask: ", this);
  }
}

export class TaskRef {
  name: string;
  dependencies?: string[];
  constructor(name: string, dependencies?: string[]) {
    this.name = name;
    this.dependencies = dependencies;
    console.log("TaskRef :", this);
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
    console.log(this);
  }
}
