TODO
- Deployment 
- Split NodeAgent into Nodecontroller and NodeAgent
  - NodeController - responsible for garbage collection of pods, collect health metrics
  - NodeAgent - responsible for pod lifecycle management
- Various deployment strategies
- Disruption budget
- Delta updates
- Concurrent updates
- StatefulSet - stable pod name
- DaemonSet - one pod per node
- Batch Job - run to completion
- More principled way to track who owns what resources / fields
- More principled way to track who is allowed to do what (change the fields)
- Proper state machine for the controller / nodeagent / allocator