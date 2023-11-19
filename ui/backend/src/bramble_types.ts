interface TaskSpec {
    image: string
    command: string[]
}

interface PLtask {
    name: string
    spec: TaskSpec
    dependencies: string[]
}

interface TaskRef {
    name: string
    dependencies: string[]
}

export interface Pipeline {
    metadata: {
        name: string
        namespace: string
    }
    spec: {
       tasks: PLtask[] 
       taskRefs: TaskRef[]
    }
}
