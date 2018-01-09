function result = shiftmiFloatEasy( x, y, varargin )
%   SHIFTMIFLOATEASY Easier interface to calculating shifted mutual in parallel.
%   SHIFTMIFLOATEASY(x, y) shifts data vectors and calculates mutual
%   information at each step.
%   x A row vector describing the first data series.
%   y A row vector describing the second data series.
%   Note: Both x and y need to have the same size.
%   Optionally the following name-value pair can be specified:
%   'shiftRange', [from to]         (default: [-500 500])
%   'binSizes',   [x-axis y-axis]   (default: [10 10])
%   'shiftSteps',  stepNumber       (default: 1)
x = single(x);
y = single(y);
p = inputParser;
addParameter(p, 'shiftRange', [-500 500]);
addParameter(p, 'binSizes', [10 10]);
addParameter(p, 'shiftSteps', 1);
parse(p, varargin{:});
minmaxOfData = [minmax(x) minmax(y)];
shiftRange = int32(p.Results.shiftRange);
binSizes = int32(p.Results.binSizes);
shiftSteps = int32(p.Results.shiftSteps);
result = shiftmiFloat(shiftRange, binSizes, minmaxOfData, x, y, shiftSteps);
end
