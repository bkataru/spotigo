window.BENCHMARK_DATA = {
  "lastUpdate": 1768415131576,
  "repoUrl": "https://github.com/bkataru/spotigo",
  "entries": {
    "Go Benchmark": [
      {
        "commit": {
          "author": {
            "email": "baalateja.k@gmail.com",
            "name": "bkataru",
            "username": "bkataru"
          },
          "committer": {
            "email": "baalateja.k@gmail.com",
            "name": "bkataru",
            "username": "bkataru"
          },
          "distinct": true,
          "id": "ecab2b1c61594cb1f6adb9684c6d3cd72db5ea23",
          "message": "fix: create gh-pages branch if missing for benchmark storage",
          "timestamp": "2026-01-14T23:54:39+05:30",
          "tree_id": "1f877baa02065f5d8ed21665e23420f77b6e09d3",
          "url": "https://github.com/bkataru/spotigo/commit/ecab2b1c61594cb1f6adb9684c6d3cd72db5ea23"
        },
        "date": 1768415131152,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkStore_Add",
            "value": 516.2,
            "unit": "ns/op\t     448 B/op\t       2 allocs/op",
            "extra": "475801 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Add - ns/op",
            "value": 516.2,
            "unit": "ns/op",
            "extra": "475801 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Add - B/op",
            "value": 448,
            "unit": "B/op",
            "extra": "475801 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Add - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "475801 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=10",
            "value": 1155,
            "unit": "ns/op\t    3752 B/op\t       7 allocs/op",
            "extra": "96484 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=10 - ns/op",
            "value": 1155,
            "unit": "ns/op",
            "extra": "96484 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=10 - B/op",
            "value": 3752,
            "unit": "B/op",
            "extra": "96484 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=10 - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "96484 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=50",
            "value": 6073,
            "unit": "ns/op\t   16744 B/op\t      11 allocs/op",
            "extra": "19710 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=50 - ns/op",
            "value": 6073,
            "unit": "ns/op",
            "extra": "19710 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=50 - B/op",
            "value": 16744,
            "unit": "B/op",
            "extra": "19710 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=50 - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "19710 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=50",
            "value": 7062,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16957 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=50 - ns/op",
            "value": 7062,
            "unit": "ns/op",
            "extra": "16957 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=50 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16957 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=50 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16957 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=100",
            "value": 14126,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8410 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=100 - ns/op",
            "value": 14126,
            "unit": "ns/op",
            "extra": "8410 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=100 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8410 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=100 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8410 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Save",
            "value": 1999717,
            "unit": "ns/op\t  427357 B/op\t     160 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Save - ns/op",
            "value": 1999717,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Save - B/op",
            "value": 427357,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Save - allocs/op",
            "value": 160,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Load",
            "value": 2219619,
            "unit": "ns/op\t  322256 B/op\t     875 allocs/op",
            "extra": "51 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Load - ns/op",
            "value": 2219619,
            "unit": "ns/op",
            "extra": "51 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Load - B/op",
            "value": 322256,
            "unit": "B/op",
            "extra": "51 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Load - allocs/op",
            "value": 875,
            "unit": "allocs/op",
            "extra": "51 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=128",
            "value": 126.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "867777 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=128 - ns/op",
            "value": 126.6,
            "unit": "ns/op",
            "extra": "867777 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=128 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "867777 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=128 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "867777 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=384",
            "value": 366,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "327319 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=384 - ns/op",
            "value": 366,
            "unit": "ns/op",
            "extra": "327319 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=384 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "327319 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=384 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "327319 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Count",
            "value": 5.806,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "20656269 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Count - ns/op",
            "value": 5.806,
            "unit": "ns/op",
            "extra": "20656269 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Count - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "20656269 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Count - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "20656269 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_CountByType",
            "value": 3527,
            "unit": "ns/op\t     256 B/op\t       2 allocs/op",
            "extra": "33380 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_CountByType - ns/op",
            "value": 3527,
            "unit": "ns/op",
            "extra": "33380 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_CountByType - B/op",
            "value": 256,
            "unit": "B/op",
            "extra": "33380 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_CountByType - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "33380 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Clear",
            "value": 460.4,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "256150 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Clear - ns/op",
            "value": 460.4,
            "unit": "ns/op",
            "extra": "256150 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Clear - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "256150 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Clear - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "256150 times\n4 procs"
          }
        ]
      }
    ]
  }
}