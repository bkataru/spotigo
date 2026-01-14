window.BENCHMARK_DATA = {
  "lastUpdate": 1768419484740,
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
      },
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
          "id": "5b1f2215216861d090c7bb984cbbc149346f3f16",
          "message": "chore: remove build artifacts and unnecessary documentation",
          "timestamp": "2026-01-15T00:13:21+05:30",
          "tree_id": "53ac756558052133222b587170eef006c426da0b",
          "url": "https://github.com/bkataru/spotigo/commit/5b1f2215216861d090c7bb984cbbc149346f3f16"
        },
        "date": 1768416246449,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkStore_Add",
            "value": 461.8,
            "unit": "ns/op\t     391 B/op\t       2 allocs/op",
            "extra": "460540 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Add - ns/op",
            "value": 461.8,
            "unit": "ns/op",
            "extra": "460540 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Add - B/op",
            "value": 391,
            "unit": "B/op",
            "extra": "460540 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Add - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "460540 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=10",
            "value": 1159,
            "unit": "ns/op\t    3752 B/op\t       7 allocs/op",
            "extra": "98949 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=10 - ns/op",
            "value": 1159,
            "unit": "ns/op",
            "extra": "98949 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=10 - B/op",
            "value": 3752,
            "unit": "B/op",
            "extra": "98949 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=10 - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "98949 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=50",
            "value": 5898,
            "unit": "ns/op\t   16744 B/op\t      11 allocs/op",
            "extra": "20332 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=50 - ns/op",
            "value": 5898,
            "unit": "ns/op",
            "extra": "20332 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=50 - B/op",
            "value": 16744,
            "unit": "B/op",
            "extra": "20332 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=50 - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "20332 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=50",
            "value": 7054,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17002 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=50 - ns/op",
            "value": 7054,
            "unit": "ns/op",
            "extra": "17002 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=50 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17002 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=50 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17002 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=100",
            "value": 14087,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8401 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=100 - ns/op",
            "value": 14087,
            "unit": "ns/op",
            "extra": "8401 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=100 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8401 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=100 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8401 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Save",
            "value": 1974574,
            "unit": "ns/op\t  459474 B/op\t     162 allocs/op",
            "extra": "51 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Save - ns/op",
            "value": 1974574,
            "unit": "ns/op",
            "extra": "51 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Save - B/op",
            "value": 459474,
            "unit": "B/op",
            "extra": "51 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Save - allocs/op",
            "value": 162,
            "unit": "allocs/op",
            "extra": "51 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Load",
            "value": 2228457,
            "unit": "ns/op\t  322268 B/op\t     875 allocs/op",
            "extra": "52 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Load - ns/op",
            "value": 2228457,
            "unit": "ns/op",
            "extra": "52 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Load - B/op",
            "value": 322268,
            "unit": "B/op",
            "extra": "52 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Load - allocs/op",
            "value": 875,
            "unit": "allocs/op",
            "extra": "52 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=128",
            "value": 126.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "948686 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=128 - ns/op",
            "value": 126.4,
            "unit": "ns/op",
            "extra": "948686 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=128 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "948686 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=128 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "948686 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=384",
            "value": 365.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "326600 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=384 - ns/op",
            "value": 365.3,
            "unit": "ns/op",
            "extra": "326600 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=384 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "326600 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=384 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "326600 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Count",
            "value": 5.763,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "20322986 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Count - ns/op",
            "value": 5.763,
            "unit": "ns/op",
            "extra": "20322986 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Count - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "20322986 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Count - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "20322986 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_CountByType",
            "value": 3578,
            "unit": "ns/op\t     256 B/op\t       2 allocs/op",
            "extra": "33265 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_CountByType - ns/op",
            "value": 3578,
            "unit": "ns/op",
            "extra": "33265 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_CountByType - B/op",
            "value": 256,
            "unit": "B/op",
            "extra": "33265 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_CountByType - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "33265 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Clear",
            "value": 405.7,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "305773 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Clear - ns/op",
            "value": 405.7,
            "unit": "ns/op",
            "extra": "305773 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Clear - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "305773 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Clear - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "305773 times\n4 procs"
          }
        ]
      },
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
          "id": "132a8d0bfa83e4f7a7283b4237ba34636ed179d4",
          "message": "Fix Windows OAuth browser opening and success page encoding\n\n- Fix openBrowser() on Windows to use rundll32 instead of cmd /c start\n  The & characters in OAuth URLs were being interpreted as command\n  separators, causing redirect_uri and other params to be truncated\n\n- Fix authentication success page character encoding\n  Added charset=utf-8 to Content-Type header and meta tag\n  Changed emoji to HTML entity for reliable rendering\n\n- Add TODO.md for tracking pending improvements",
          "timestamp": "2026-01-15T01:07:18+05:30",
          "tree_id": "ec6ae3e84cd42e43f11b37361db3ad86f85f7cc3",
          "url": "https://github.com/bkataru/spotigo/commit/132a8d0bfa83e4f7a7283b4237ba34636ed179d4"
        },
        "date": 1768419483921,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkStore_Add",
            "value": 432.9,
            "unit": "ns/op\t     364 B/op\t       2 allocs/op",
            "extra": "320682 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Add - ns/op",
            "value": 432.9,
            "unit": "ns/op",
            "extra": "320682 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Add - B/op",
            "value": 364,
            "unit": "B/op",
            "extra": "320682 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Add - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "320682 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=10",
            "value": 1183,
            "unit": "ns/op\t    3752 B/op\t       7 allocs/op",
            "extra": "100486 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=10 - ns/op",
            "value": 1183,
            "unit": "ns/op",
            "extra": "100486 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=10 - B/op",
            "value": 3752,
            "unit": "B/op",
            "extra": "100486 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=10 - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "100486 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=50",
            "value": 5970,
            "unit": "ns/op\t   16744 B/op\t      11 allocs/op",
            "extra": "19977 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=50 - ns/op",
            "value": 5970,
            "unit": "ns/op",
            "extra": "19977 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=50 - B/op",
            "value": 16744,
            "unit": "B/op",
            "extra": "19977 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_AddBatch/size=50 - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "19977 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=50",
            "value": 7033,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16861 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=50 - ns/op",
            "value": 7033,
            "unit": "ns/op",
            "extra": "16861 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=50 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16861 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=50 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16861 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=100",
            "value": 14092,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8556 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=100 - ns/op",
            "value": 14092,
            "unit": "ns/op",
            "extra": "8556 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=100 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8556 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SearchSimilarity/size=100 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8556 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Save",
            "value": 2015269,
            "unit": "ns/op\t  438896 B/op\t     161 allocs/op",
            "extra": "51 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Save - ns/op",
            "value": 2015269,
            "unit": "ns/op",
            "extra": "51 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Save - B/op",
            "value": 438896,
            "unit": "B/op",
            "extra": "51 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Save - allocs/op",
            "value": 161,
            "unit": "allocs/op",
            "extra": "51 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Load",
            "value": 2226400,
            "unit": "ns/op\t  322268 B/op\t     875 allocs/op",
            "extra": "54 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Load - ns/op",
            "value": 2226400,
            "unit": "ns/op",
            "extra": "54 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Load - B/op",
            "value": 322268,
            "unit": "B/op",
            "extra": "54 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_SaveLoad/Load - allocs/op",
            "value": 875,
            "unit": "allocs/op",
            "extra": "54 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=128",
            "value": 126.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "899210 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=128 - ns/op",
            "value": 126.7,
            "unit": "ns/op",
            "extra": "899210 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=128 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "899210 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=128 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "899210 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=384",
            "value": 366.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "327476 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=384 - ns/op",
            "value": 366.4,
            "unit": "ns/op",
            "extra": "327476 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=384 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "327476 times\n4 procs"
          },
          {
            "name": "BenchmarkCosineSimilarity/dim=384 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "327476 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Count",
            "value": 5.767,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "20537373 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Count - ns/op",
            "value": 5.767,
            "unit": "ns/op",
            "extra": "20537373 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Count - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "20537373 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Count - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "20537373 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_CountByType",
            "value": 3520,
            "unit": "ns/op\t     256 B/op\t       2 allocs/op",
            "extra": "33547 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_CountByType - ns/op",
            "value": 3520,
            "unit": "ns/op",
            "extra": "33547 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_CountByType - B/op",
            "value": 256,
            "unit": "B/op",
            "extra": "33547 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_CountByType - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "33547 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Clear",
            "value": 422.8,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "270214 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Clear - ns/op",
            "value": 422.8,
            "unit": "ns/op",
            "extra": "270214 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Clear - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "270214 times\n4 procs"
          },
          {
            "name": "BenchmarkStore_Clear - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "270214 times\n4 procs"
          }
        ]
      }
    ]
  }
}