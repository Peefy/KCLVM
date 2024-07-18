func @cal_grid_from_dispacth_workloads(%arg1: index, %arg2: index, %arg3: index) {
  %c1 = arith.constant 1 : index
  %0 = affine.apply affine_map<()[s0] -> (s0 ceildiv 64)>()[%arg1]
  %1 = affine.apply affine_map<()[s0] -> (s0 ceildiv 64)>()[%arg2]
  return %arg1, %0, %c1 : index, index, index
}