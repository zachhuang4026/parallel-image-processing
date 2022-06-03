# Parallel Image Processing
This is a project from MPCS-52060 Parallel Programming course. The project speeds up the image processing progress by parallelism coded in Golang. This project includes 4 version of parallelism techniques, and will be demonstrated with its speed up graph.
- Pipeline by Channels
- Bulk Synchronous Parallel (BSP) 
- Work Stealing Algorithms
- Work Balancing Algorithms

The parallel processing is to divide one image into chucks and apply to effect at the same time.
<p align = 'center'>
<img src = 'https://github.com/zachhuang4026/parallel-image-processing/blob/main/demo.png' width="600">
</p>

# Speed Up Graph
The input images all include 3 version based on its resolution (small/mixture/big). 

### Pipeline 
<p align = 'center'>
<img src = 'https://github.com/zachhuang4026/parallel-image-processing/blob/main/speedup_pipeline.png' width="600">
</p>

### BSP
<p align = 'center'>
<img src = 'https://github.com/zachhuang4026/parallel-image-processing/blob/main/speedup_bsp.png' width="600">
</p>

### Work Stealing
<p align = 'center'>
<img src = 'https://github.com/zachhuang4026/parallel-image-processing/blob/main/speedup_steal.png' width="600">
</p>

### Work Balancing
<p align = 'center'>
<img src = 'https://github.com/zachhuang4026/parallel-image-processing/blob/main/speedup_balance.png' width="600">
</p>
