package com.ndisk.webdav;

public interface FileSystem {
    void Mkdir(String name) throws Exception;

    File OpenFile(String name) throws Exception;

    void RemoveAll(String name) throws Exception;

    void Rename(String oldName, String newName) throws Exception;

    java.io.File Stat(String name) throws Exception;
}
