# Bramble - A Kubernetes Native CI/CD Framework

## Developer Experience
- Pipelines as Kubernetes manifests
- Pipelines can have handmade tasks or plug and play pre-applied tasks

### Demo ideas
- Pipelines to build, test, and deploy different components of Bramble
- C program to generate pngs, merging to master posts the generated image to instagram

### Interim targets
- Nov 25th: Operator working against a test repo. Rendering execution
- Nov 30th: CSS and other goodness
- Dec 5th: Report completed

**Note: Complexity of tasks will be measured in terms of video game difficulties.**
1. I'm too young to die: Minor bug fixes, typos, code quality (Measured in minutes)
2. Hey, not too rough: Implmentation of functions in a feature, rewriting certian segments.
3. Hurt me plenty: Implemetation of individual features within a component.
4. Ultra violence: Interop between components, rewriting large portions of code, or changes to the stack.
5. Ultrakill: Entire components, architectural decisions.

### Current Tasks
#### Operator
- Have the operator create an in-memory database of tasks and pipelines (3)
- Allow the operator to manage execution of pipelines. (4)
- This may involve the creation of another CRD. (5)
- New CRD would be an execution and would let the operator provision resources such as containers, volumes, and more. (4)
- Allow for nested pipelines (3)
- Implement staging to stagger out different parts of the pipelines (3)

#### UI
- Render the execution of a test pipeline (3)
- Allow users to click into individual steps and view logs (4)
- Show currently applied pipelines and task in the cluster (3)

#### Git Provider App
- The complexity of this project seems to require a GitHub app to handle interop between GitHub and the Operator (5)
- Establish a method of listening for GitHub events such as PRS, Merges, Etc (4)
- Create an easy YAML schema to show what events to look out for (CRD or other method) (4)
- Allow git provider app to emit executions for operator to handle. (3)
