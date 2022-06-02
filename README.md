# Parallel Image Processing
This is a project from MPCS-52060 Parallel Programming course. The project speeds up the image processing progress by parallelism coded in Golang. This project includes 4 version of parallelism techniques, and will be demonstrated with its speed up graph.
- Pipeline by Channels
- Bulk Synchronous Parallel (BSP) 
- Work Stealing Algorithms
- Work Balancing Algorithms

# Speed Up Graph
### Pipeline 
<p align = 'center'>
<img src = 'https://github.com/zachhuang4026/parallel-image-processing/blob/main/speedup_pipeline.png' width="300">
</p>
![alt text](https://github.com/zachhuang4026/parallel-image-processing/blob/main/speedup_pipeline.pdf)

### BSP
<p align = 'center'>
<img src = 'https://github.com/zachhuang4026/parallel-image-processing/blob/main/speedup_bsp.pdf' width="300">
</p>

### Work Stealing
<p align = 'center'>
<img src = 'https://github.com/zachhuang4026/parallel-image-processing/blob/main/speedup_steal.pdf' width="300">
</p>

### Work Balancing
<p align = 'center'>
<img src = 'https://github.com/zachhuang4026/parallel-image-processing/blob/main/speedup_balance.pdf' width="300">
</p>
