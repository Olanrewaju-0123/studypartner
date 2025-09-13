"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { ArrowLeft, CheckCircle, AlertCircle } from "lucide-react";
import Link from "next/link";
import FileUpload from "@/components/FileUpload";
import { Note } from "@/types";
import ThemeToggle from "@/components/ThemeToggle";

export default function UploadPage() {
  const router = useRouter();
  const [uploadedNote, setUploadedNote] = useState<Note | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleUploadSuccess = (note: Note) => {
    setUploadedNote(note);
    setError(null);
  };

  const handleUploadError = (errorMessage: string) => {
    setError(errorMessage);
    setUploadedNote(null);
  };

  const goToStudy = () => {
    if (uploadedNote) {
      router.push(`/study/${uploadedNote.id}`);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      {/* Header */}
      <div className="bg-white dark:bg-gray-800 shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center">
              <Link
                href="/"
                className="flex items-center text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white text-sm sm:text-base"
              >
                <ArrowLeft className="h-4 w-4 sm:h-5 sm:w-5 mr-1 sm:mr-2" />
                <span className="hidden sm:inline">Back to Home</span>
                <span className="sm:hidden">Back</span>
              </Link>
            </div>
            <h1 className="text-lg sm:text-xl font-semibold text-gray-900 dark:text-white">
              Upload Notes
            </h1>
            <div className="flex items-center">
              <ThemeToggle />
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12">
        <div className="text-center mb-8 sm:mb-12">
          <h2 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white mb-4">
            Upload Your Study Materials
          </h2>
          <p className="text-base sm:text-lg text-gray-600 dark:text-gray-300 max-w-2xl mx-auto px-4">
            Upload your PDF, DOCX, or TXT files to generate AI-powered
            summaries, flashcards, and quizzes.
          </p>
        </div>

        {/* Upload Component */}
        <div className="mb-8">
          <FileUpload
            onUploadSuccess={handleUploadSuccess}
            onUploadError={handleUploadError}
          />
        </div>

        {/* Success Message */}
        {uploadedNote && (
          <div className="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-6 mb-8">
            <div className="flex items-center">
              <CheckCircle className="h-6 w-6 text-green-500 dark:text-green-400 mr-3" />
              <div>
                <h3 className="text-lg font-semibold text-green-900 dark:text-green-100">
                  Upload Successful!
                </h3>
                <p className="text-green-700 dark:text-green-300">
                  Your file "{uploadedNote.file_name}" has been processed
                  successfully.
                </p>
              </div>
            </div>
            <div className="mt-4">
              <button
                onClick={goToStudy}
                className="bg-green-600 dark:bg-green-500 text-white px-6 py-2 rounded-lg hover:bg-green-700 dark:hover:bg-green-600 transition-colors"
              >
                Start Studying
              </button>
            </div>
          </div>
        )}

        {/* Error Message */}
        {error && (
          <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6 mb-8">
            <div className="flex items-center">
              <AlertCircle className="h-6 w-6 text-red-500 dark:text-red-400 mr-3" />
              <div>
                <h3 className="text-lg font-semibold text-red-900 dark:text-red-100">
                  Upload Failed
                </h3>
                <p className="text-red-700 dark:text-red-300">{error}</p>
              </div>
            </div>
          </div>
        )}

        {/* Features Preview */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 mt-8 sm:mt-12">
          <div className="bg-white dark:bg-gray-800 p-4 sm:p-6 rounded-lg shadow-sm hover:shadow-lg hover:-translate-y-1 transition-all duration-300 ease-in-out transform">
            <div className="bg-blue-100 dark:bg-blue-900/30 w-10 h-10 sm:w-12 sm:h-12 rounded-lg flex items-center justify-center mb-3 sm:mb-4 transition-transform duration-300 hover:scale-110">
              <span className="text-blue-600 dark:text-blue-400 font-bold text-lg sm:text-xl">📝</span>
            </div>
            <h3 className="text-base sm:text-lg font-semibold text-gray-900 dark:text-white mb-2 transition-colors duration-300">
              Smart Summaries
            </h3>
            <p className="text-gray-600 dark:text-gray-300 text-xs sm:text-sm transition-colors duration-300">
              Get concise summaries that highlight key concepts and main points
              from your notes.
            </p>
          </div>

          <div className="bg-white dark:bg-gray-800 p-4 sm:p-6 rounded-lg shadow-sm hover:shadow-lg hover:-translate-y-1 transition-all duration-300 ease-in-out transform">
            <div className="bg-green-100 dark:bg-green-900/30 w-10 h-10 sm:w-12 sm:h-12 rounded-lg flex items-center justify-center mb-3 sm:mb-4 transition-transform duration-300 hover:scale-110">
              <span className="text-green-600 dark:text-green-400 font-bold text-lg sm:text-xl">🎯</span>
            </div>
            <h3 className="text-base sm:text-lg font-semibold text-gray-900 dark:text-white mb-2 transition-colors duration-300">
              Flashcards
            </h3>
            <p className="text-gray-600 dark:text-gray-300 text-xs sm:text-sm transition-colors duration-300">
              Generate interactive flashcards with questions and answers to test
              your knowledge.
            </p>
          </div>

          <div className="bg-white dark:bg-gray-800 p-4 sm:p-6 rounded-lg shadow-sm hover:shadow-lg hover:-translate-y-1 transition-all duration-300 ease-in-out transform sm:col-span-2 lg:col-span-1">
            <div className="bg-purple-100 dark:bg-purple-900/30 w-10 h-10 sm:w-12 sm:h-12 rounded-lg flex items-center justify-center mb-3 sm:mb-4 transition-transform duration-300 hover:scale-110">
              <span className="text-purple-600 dark:text-purple-400 font-bold text-lg sm:text-xl">❓</span>
            </div>
            <h3 className="text-base sm:text-lg font-semibold text-gray-900 dark:text-white mb-2 transition-colors duration-300">
              Quizzes
            </h3>
            <p className="text-gray-600 dark:text-gray-300 text-xs sm:text-sm transition-colors duration-300">
              Create multiple choice quizzes to assess your understanding of the
              material.
            </p>
          </div>
        </div>

        {/* Supported Formats */}
        <div className="mt-8 sm:mt-12 bg-white dark:bg-gray-800 rounded-lg shadow-sm p-4 sm:p-6">
          <h3 className="text-base sm:text-lg font-semibold text-gray-900 dark:text-white mb-4">
            Supported File Formats
          </h3>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            <div className="flex items-center">
              <div className="bg-red-100 dark:bg-red-900/30 w-8 h-8 sm:w-10 sm:h-10 rounded-lg flex items-center justify-center mr-3">
                <span className="text-red-600 dark:text-red-400 font-bold text-xs sm:text-sm">PDF</span>
              </div>
              <div>
                <p className="font-medium text-gray-900 dark:text-white text-sm sm:text-base">PDF Documents</p>
                <p className="text-xs sm:text-sm text-gray-600 dark:text-gray-300">
                  Text extraction from PDF files
                </p>
              </div>
            </div>
            <div className="flex items-center">
              <div className="bg-blue-100 dark:bg-blue-900/30 w-8 h-8 sm:w-10 sm:h-10 rounded-lg flex items-center justify-center mr-3">
                <span className="text-blue-600 dark:text-blue-400 font-bold text-xs sm:text-sm">DOCX</span>
              </div>
              <div>
                <p className="font-medium text-gray-900 dark:text-white text-sm sm:text-base">Word Documents</p>
                <p className="text-xs sm:text-sm text-gray-600 dark:text-gray-300">Microsoft Word files</p>
              </div>
            </div>
            <div className="flex items-center sm:col-span-2 lg:col-span-1">
              <div className="bg-gray-100 dark:bg-gray-700 w-8 h-8 sm:w-10 sm:h-10 rounded-lg flex items-center justify-center mr-3">
                <span className="text-gray-600 dark:text-gray-300 font-bold text-xs sm:text-sm">TXT</span>
              </div>
              <div>
                <p className="font-medium text-gray-900 dark:text-white text-sm sm:text-base">Text Files</p>
                <p className="text-xs sm:text-sm text-gray-600 dark:text-gray-300">Plain text documents</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
