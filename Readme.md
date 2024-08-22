TODO for prototype
- Watchdog
- Deployment 
- Split NodeAgent into Nodecontroller and NodeAgent
  - NodeController - responsible for garbage collection of pods, collect health metrics
  - NodeAgent - responsible for pod lifecycle management
- Various deployment strategies
- Disruption budget
- Delta updates
- Concurrent updates
- StatefulSet - stable pod name
- More principled way to track state changes -- who owns what resources / fields / objects
- Proper state machine for the controller / nodeagent / allocator


