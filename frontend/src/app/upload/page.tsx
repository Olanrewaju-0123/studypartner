"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { ArrowLeft, CheckCircle, AlertCircle } from "lucide-react";
import Link from "next/link";
import FileUpload from "@/components/FileUpload";
import { Note } from "@/types";

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
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center">
              <Link
                href="/"
                className="flex items-center text-gray-600 hover:text-gray-900"
              >
                <ArrowLeft className="h-5 w-5 mr-2" />
                Back to Home
              </Link>
            </div>
            <h1 className="text-xl font-semibold text-gray-900">
              Upload Notes
            </h1>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="text-center mb-12">
          <h2 className="text-3xl font-bold text-gray-900 mb-4">
            Upload Your Study Materials
          </h2>
          <p className="text-lg text-gray-600 max-w-2xl mx-auto">
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
          <div className="bg-green-50 border border-green-200 rounded-lg p-6 mb-8">
            <div className="flex items-center">
              <CheckCircle className="h-6 w-6 text-green-500 mr-3" />
              <div>
                <h3 className="text-lg font-semibold text-green-900">
                  Upload Successful!
                </h3>
                <p className="text-green-700">
                  Your file "{uploadedNote.file_name}" has been processed
                  successfully.
                </p>
              </div>
            </div>
            <div className="mt-4">
              <button
                onClick={goToStudy}
                className="bg-green-600 text-white px-6 py-2 rounded-lg hover:bg-green-700 transition-colors"
              >
                Start Studying
              </button>
            </div>
          </div>
        )}

        {/* Error Message */}
        {error && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-6 mb-8">
            <div className="flex items-center">
              <AlertCircle className="h-6 w-6 text-red-500 mr-3" />
              <div>
                <h3 className="text-lg font-semibold text-red-900">
                  Upload Failed
                </h3>
                <p className="text-red-700">{error}</p>
              </div>
            </div>
          </div>
        )}

        {/* Features Preview */}
        <div className="grid md:grid-cols-3 gap-6 mt-12">
          <div className="bg-white p-6 rounded-lg shadow-sm">
            <div className="bg-blue-100 w-12 h-12 rounded-lg flex items-center justify-center mb-4">
              <span className="text-blue-600 font-bold">📝</span>
            </div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">
              Smart Summaries
            </h3>
            <p className="text-gray-600 text-sm">
              Get concise summaries that highlight key concepts and main points
              from your notes.
            </p>
          </div>

          <div className="bg-white p-6 rounded-lg shadow-sm">
            <div className="bg-green-100 w-12 h-12 rounded-lg flex items-center justify-center mb-4">
              <span className="text-green-600 font-bold">🎯</span>
            </div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">
              Flashcards
            </h3>
            <p className="text-gray-600 text-sm">
              Generate interactive flashcards with questions and answers to test
              your knowledge.
            </p>
          </div>

          <div className="bg-white p-6 rounded-lg shadow-sm">
            <div className="bg-purple-100 w-12 h-12 rounded-lg flex items-center justify-center mb-4">
              <span className="text-purple-600 font-bold">❓</span>
            </div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">
              Quizzes
            </h3>
            <p className="text-gray-600 text-sm">
              Create multiple choice quizzes to assess your understanding of the
              material.
            </p>
          </div>
        </div>

        {/* Supported Formats */}
        <div className="mt-12 bg-white rounded-lg shadow-sm p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">
            Supported File Formats
          </h3>
          <div className="grid md:grid-cols-3 gap-4">
            <div className="flex items-center">
              <div className="bg-red-100 w-10 h-10 rounded-lg flex items-center justify-center mr-3">
                <span className="text-red-600 font-bold text-sm">PDF</span>
              </div>
              <div>
                <p className="font-medium text-gray-900">PDF Documents</p>
                <p className="text-sm text-gray-600">
                  Text extraction from PDF files
                </p>
              </div>
            </div>
            <div className="flex items-center">
              <div className="bg-blue-100 w-10 h-10 rounded-lg flex items-center justify-center mr-3">
                <span className="text-blue-600 font-bold text-sm">DOCX</span>
              </div>
              <div>
                <p className="font-medium text-gray-900">Word Documents</p>
                <p className="text-sm text-gray-600">Microsoft Word files</p>
              </div>
            </div>
            <div className="flex items-center">
              <div className="bg-gray-100 w-10 h-10 rounded-lg flex items-center justify-center mr-3">
                <span className="text-gray-600 font-bold text-sm">TXT</span>
              </div>
              <div>
                <p className="font-medium text-gray-900">Text Files</p>
                <p className="text-sm text-gray-600">Plain text documents</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
