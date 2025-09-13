"use client";

import { useState, useCallback } from "react";
import { useDropzone } from "react-dropzone";
import { Upload, File, X, Loader2 } from "lucide-react";
import { apiClient } from "@/utils/api";
import { Note } from "@/types";

interface FileUploadProps {
  onUploadSuccess: (note: Note) => void;
  onUploadError: (error: string) => void;
}

export default function FileUpload({
  onUploadSuccess,
  onUploadError,
}: FileUploadProps) {
  const [uploading, setUploading] = useState(false);
  const [uploadedFile, setUploadedFile] = useState<File | null>(null);

  const onDrop = useCallback(
    async (acceptedFiles: File[]) => {
      const file = acceptedFiles[0];
      if (!file) return;

      setUploadedFile(file);
      setUploading(true);

      try {
        // Convert file to base64
        const base64 = await fileToBase64(file);

        // Upload to backend
        const note = await apiClient.uploadNote({
          file: base64,
          name: file.name,
        });

        onUploadSuccess(note);
        setUploadedFile(null);
      } catch (error) {
        onUploadError(error instanceof Error ? error.message : "Upload failed");
      } finally {
        setUploading(false);
      }
    },
    [onUploadSuccess, onUploadError]
  );

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      "application/pdf": [".pdf"],
      "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
        [".docx"],
      "text/plain": [".txt"],
    },
    maxFiles: 1,
    disabled: uploading,
  });

  const fileToBase64 = (file: File): Promise<string> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => {
        const result = reader.result as string;
        // Remove data URL prefix
        const base64 = result.split(",")[1];
        resolve(base64);
      };
      reader.onerror = (error) => reject(error);
    });
  };

  const removeFile = () => {
    setUploadedFile(null);
  };

  return (
    <div className="w-full max-w-2xl mx-auto">
      <div
        {...getRootProps()}
        className={`
          border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors
          ${
            isDragActive
              ? "border-blue-500 dark:border-blue-400 bg-blue-50 dark:bg-blue-900/20"
              : "border-gray-300 dark:border-gray-600 hover:border-gray-400 dark:hover:border-gray-500"
          }
          ${uploading ? "opacity-50 cursor-not-allowed" : ""}
        `}
      >
        <input {...getInputProps()} />

        {uploading ? (
          <div className="flex flex-col items-center">
            <Loader2 className="h-12 w-12 text-blue-500 dark:text-blue-400 animate-spin mb-4" />
            <p className="text-lg font-medium text-gray-700 dark:text-gray-300">Uploading...</p>
            <p className="text-sm text-gray-500 dark:text-gray-400">Processing your file</p>
          </div>
        ) : uploadedFile ? (
          <div className="flex items-center justify-center space-x-4">
            <File className="h-8 w-8 text-green-500 dark:text-green-400" />
            <div className="text-left">
              <p className="font-medium text-gray-900 dark:text-white">{uploadedFile.name}</p>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                {(uploadedFile.size / 1024 / 1024).toFixed(2)} MB
              </p>
            </div>
            <button
              onClick={removeFile}
              className="p-1 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full"
            >
              <X className="h-4 w-4 text-gray-500 dark:text-gray-400" />
            </button>
          </div>
        ) : (
          <div className="flex flex-col items-center">
            <Upload className="h-12 w-12 text-gray-400 dark:text-gray-500 mb-4" />
            <p className="text-lg font-medium text-gray-700 dark:text-gray-300 mb-2">
              {isDragActive ? "Drop your file here" : "Upload your study notes"}
            </p>
            <p className="text-sm text-gray-500 dark:text-gray-400 mb-4">
              Drag and drop or click to select
            </p>
            <p className="text-xs text-gray-400 dark:text-gray-500">
              Supports PDF, DOCX, and TXT files
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
