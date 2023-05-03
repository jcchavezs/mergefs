# mergefs

MergeFS allows to merge different filesystems into a single one, trying to keep the ancillary interfaces (`GlobFS`, `FileReaderFS`) to avoid side effects when calling such methods in the merged FS. Standard methods for reading a file or searching globs might take different paths depending on whether the filesystem implements one ancillary interface or not ending up in slightly different results. By keeping the public API of the filesystems we aim to avoid such differences.

**Important:** Clashing paths among merging FS have been skipped on purpose to keep this library simple by respecting the ordering of the FS being merged, meaning that when looking for a file, the search ends as soon as the file is found.

## Getting started

```go
fsA := fstest.MapFS{"a.txt": &fstest.MapFile{Data: []byte("test a")}}
fsB := fstest.MapFS{"b.txt": &fstest.MapFile{Data: []byte("text b")}}

mfs := mergefs.Merge(fsA, fsB)
if _, err := mfs.Open("a.txt"); err != nil {
  // ...
}
```

## Related work

`MergeFS` is highy inspired by https://github.com/laher/mergefs and https://github.com/yalue/merged_fs.
