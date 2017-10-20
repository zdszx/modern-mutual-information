/**
 * Copyright 2017 Thomas Leyh <leyht@informatik.uni-freiburg.de>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include <catch.hpp>
#include <cmath>
#include <algorithm>
#include "../src/utilities.h"

TEST_CASE( "Test calculation of indices in 1 dimension.", "[calculate_indices_1d]" )
{
	std::vector<float> input(1000);
	for (int i = 0; i < 1000; ++i)
	{
		input[i] = float(i) - 500.f;
	}
	auto indices = calculate_indices_1d(10, -500.f, 499.f, input.begin(), input.end());
	REQUIRE( indices.size() == 1000 );
	CHECK( indices[0] == 0 );
	CHECK( indices[23] == 0 );
	CHECK( indices[99] == 0 );
	CHECK( indices[100] == 1 );
	CHECK( indices[199] == 1 );
	CHECK( indices[990] == 9 );
	CHECK( indices[999] == 9 );
}

TEST_CASE( "Test calculation of indices in 2 dimensions.", "[calculate_indices_2d]" )
{
	std::vector<float> inputX(800);
	for (int i = 0; i < 800; ++i)
	{
		inputX[i] = float(i) - 500.f;
	}
	std::vector<float> inputY(800);
	for (int i = 0; i < 800; ++i)
	{
		inputY[i] = float(i) - 400.f;
	}
	auto indices = calculate_indices_2d(10, 10,
			inputX.front(), inputX.back(), inputY.front(), inputY.back(),
			inputX.begin(), inputX.end(), inputY.begin(), inputY.end());
	REQUIRE( indices.size() == 800 );
	CHECK( indices[0].first == 0 );
	CHECK( indices[0].second == 0 );
	CHECK( indices[79].first == 0 );
	CHECK( indices[79].second == 0 );
	CHECK( indices[80].first == 1 );
	CHECK( indices[80].second == 1 );
	CHECK( indices[799].first == 9 );
	CHECK( indices[799].second == 9 );
}

TEST_CASE( "Shifted mutual information on sinoid data." "[shifted_mutual_information]")
{
	// Fill a vector with sinoid data.
	std::vector<float> data(1000);
	float value = 0;
	for (auto& d : data)
	{
		d = std::sin(value);
		value += 0.01f;
	}
	auto result = shifted_mutual_information(-100, 101, 10, 10, -1.f, 1.f, -1.f, 1.f,
											 data.begin(), data.end(), data.begin(), data.end());
	REQUIRE( result.size() == 201 );
	auto maximum = std::max_element(result.begin(), result.end());
	REQUIRE( std::distance(result.begin(), maximum) == 100 ); // Should be exactly the middle.
	// And now check for symmetries.
	auto rbegin = result.rbegin();
	for (auto begin = result.begin(); begin != maximum; ++begin)
	{
		CHECK( *begin == Approx(*rbegin) );
		++rbegin;
	}

	auto result_s3 = shifted_mutual_information(-100, 101, 10, 10, -1.f, 1.f, -1.f, 1.f,
												data.begin(), data.end(), data.begin(), data.end(), 3);
	REQUIRE( result_s3.size() == 67 );
	auto maximum_s3 = std::max_element(result_s3.begin(), result_s3.end());
	REQUIRE( std::distance(result_s3.begin(), maximum_s3) == 33 );
}
