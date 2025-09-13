"use client";

import { useState } from "react";
import { Upload, BookOpen, Brain, Target, Search } from "lucide-react";
import Link from "next/link";
import ThemeToggle from "@/components/ThemeToggle";

export default function Home() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800">
      {/* Navigation */}
      <nav className="bg-white dark:bg-gray-800 shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <BookOpen className="h-6 w-6 sm:h-8 sm:w-8 text-blue-600 dark:text-blue-400" />
              <span className="ml-2 text-lg sm:text-xl font-bold text-gray-900 dark:text-white">
                AI Study Partner
              </span>
            </div>
            <div className="flex items-center space-x-2 sm:space-x-4">
              <ThemeToggle />
              {isAuthenticated ? (
                <>
                  <Link
                    href="/dashboard"
                    className="hidden sm:block text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 text-sm sm:text-base"
                  >
                    Dashboard
                  </Link>
                  <button className="bg-blue-600 dark:bg-blue-500 text-white px-2 sm:px-4 py-1 sm:py-2 rounded-md hover:bg-blue-700 dark:hover:bg-blue-600 text-sm sm:text-base">
                    Logout
                  </button>
                </>
              ) : (
                <>
                  <Link
                    href="/login"
                    className="hidden sm:block text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 text-sm sm:text-base"
                  >
                    Login
                  </Link>
                  <Link
                    href="/register"
                    className="bg-blue-600 dark:bg-blue-500 text-white px-2 sm:px-4 py-1 sm:py-2 rounded-md hover:bg-blue-700 dark:hover:bg-blue-600 text-sm sm:text-base"
                  >
                    <span className="hidden sm:inline">Sign Up</span>
                    <span className="sm:hidden">Sign Up</span>
                  </Link>
                </>
              )}
            </div>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 sm:py-16 lg:py-20">
        <div className="text-center">
          <h1 className="text-3xl sm:text-4xl md:text-5xl lg:text-6xl font-bold text-gray-900 dark:text-white mb-4 sm:mb-6">
            Study Smarter with
            <span className="text-blue-600 dark:text-blue-400"> AI</span>
          </h1>
          <p className="text-lg sm:text-xl text-gray-600 dark:text-gray-300 mb-6 sm:mb-8 max-w-3xl mx-auto px-4">
            Transform your notes into powerful study materials. Upload your
            documents and get AI-generated summaries, flashcards, and quizzes
            instantly.
          </p>
          <div className="flex flex-col sm:flex-row gap-3 sm:gap-4 justify-center px-4">
            <Link
              href="/upload"
              className="bg-blue-600 dark:bg-blue-500 text-white px-6 sm:px-8 py-3 rounded-lg text-base sm:text-lg font-semibold hover:bg-blue-700 dark:hover:bg-blue-600 transition-colors"
            >
              <Upload className="inline-block mr-2 h-4 w-4 sm:h-5 sm:w-5" />
              Upload Notes
            </Link>
            <Link
              href="/demo"
              className="border border-blue-600 dark:border-blue-400 text-blue-600 dark:text-blue-400 px-6 sm:px-8 py-3 rounded-lg text-base sm:text-lg font-semibold hover:bg-blue-50 dark:hover:bg-blue-900/20 transition-colors"
            >
              Try Demo
            </Link>
          </div>
        </div>
      </div>

      {/* Features Section */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 sm:py-16 lg:py-20">
        <div className="text-center mb-8 sm:mb-12 lg:mb-16">
          <h2 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white mb-4">
            Powerful AI Features
          </h2>
          <p className="text-base sm:text-lg text-gray-600 dark:text-gray-300">
            Everything you need to study effectively
          </p>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 sm:gap-8">
          <div className="bg-white dark:bg-gray-800 p-6 sm:p-8 rounded-xl shadow-lg hover:shadow-xl hover:-translate-y-1 transition-all duration-300 ease-in-out transform">
            <div className="bg-blue-100 dark:bg-blue-900/30 w-10 h-10 sm:w-12 sm:h-12 rounded-lg flex items-center justify-center mb-4 sm:mb-6 transition-transform duration-300 hover:scale-110">
              <Brain className="h-5 w-5 sm:h-6 sm:w-6 text-blue-600 dark:text-blue-400" />
            </div>
            <h3 className="text-lg sm:text-xl font-semibold text-gray-900 dark:text-white mb-3 sm:mb-4 transition-colors duration-300">
              Smart Summaries
            </h3>
            <p className="text-sm sm:text-base text-gray-600 dark:text-gray-300 transition-colors duration-300">
              Get concise, comprehensive summaries of your notes that highlight
              key concepts and main points.
            </p>
          </div>

          <div className="bg-white dark:bg-gray-800 p-6 sm:p-8 rounded-xl shadow-lg hover:shadow-xl hover:-translate-y-1 transition-all duration-300 ease-in-out transform">
            <div className="bg-green-100 dark:bg-green-900/30 w-10 h-10 sm:w-12 sm:h-12 rounded-lg flex items-center justify-center mb-4 sm:mb-6 transition-transform duration-300 hover:scale-110">
              <Target className="h-5 w-5 sm:h-6 sm:w-6 text-green-600 dark:text-green-400" />
            </div>
            <h3 className="text-lg sm:text-xl font-semibold text-gray-900 dark:text-white mb-3 sm:mb-4 transition-colors duration-300">
              Flashcards
            </h3>
            <p className="text-sm sm:text-base text-gray-600 dark:text-gray-300 transition-colors duration-300">
              Generate interactive flashcards with questions and answers to test
              your knowledge and improve retention.
            </p>
          </div>

          <div className="bg-white dark:bg-gray-800 p-6 sm:p-8 rounded-xl shadow-lg hover:shadow-xl hover:-translate-y-1 transition-all duration-300 ease-in-out transform sm:col-span-2 lg:col-span-1">
            <div className="bg-purple-100 dark:bg-purple-900/30 w-10 h-10 sm:w-12 sm:h-12 rounded-lg flex items-center justify-center mb-4 sm:mb-6 transition-transform duration-300 hover:scale-110">
              <Search className="h-5 w-5 sm:h-6 sm:w-6 text-purple-600 dark:text-purple-400" />
            </div>
            <h3 className="text-lg sm:text-xl font-semibold text-gray-900 dark:text-white mb-3 sm:mb-4 transition-colors duration-300">
              Smart Search
            </h3>
            <p className="text-sm sm:text-base text-gray-600 dark:text-gray-300 transition-colors duration-300">
              Find relevant information across all your notes using semantic
              search powered by AI embeddings.
            </p>
          </div>
        </div>
      </div>

      {/* How it Works */}
      <div className="bg-white dark:bg-gray-800 py-12 sm:py-16 lg:py-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-8 sm:mb-12 lg:mb-16">
            <h2 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white mb-4">
              How It Works
            </h2>
            <p className="text-base sm:text-lg text-gray-600 dark:text-gray-300">
              Simple steps to transform your study routine
            </p>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 sm:gap-8">
            <div className="text-center group hover:scale-105 transition-all duration-300 ease-in-out">
              <div className="bg-blue-600 dark:bg-blue-500 text-white w-10 h-10 sm:w-12 sm:h-12 rounded-full flex items-center justify-center mx-auto mb-3 sm:mb-4 text-lg sm:text-xl font-bold transition-all duration-300 group-hover:scale-110 group-hover:shadow-lg">
                1
              </div>
              <h3 className="text-base sm:text-lg font-semibold text-gray-900 dark:text-white mb-2 transition-colors duration-300">
                Upload
              </h3>
              <p className="text-sm sm:text-base text-gray-600 dark:text-gray-300 transition-colors duration-300">
                Upload your PDF, DOCX, or TXT files
              </p>
            </div>

            <div className="text-center group hover:scale-105 transition-all duration-300 ease-in-out">
              <div className="bg-blue-600 dark:bg-blue-500 text-white w-10 h-10 sm:w-12 sm:h-12 rounded-full flex items-center justify-center mx-auto mb-3 sm:mb-4 text-lg sm:text-xl font-bold transition-all duration-300 group-hover:scale-110 group-hover:shadow-lg">
                2
              </div>
              <h3 className="text-base sm:text-lg font-semibold text-gray-900 dark:text-white mb-2 transition-colors duration-300">
                Process
              </h3>
              <p className="text-sm sm:text-base text-gray-600 dark:text-gray-300 transition-colors duration-300">
                AI extracts and analyzes your content
              </p>
            </div>

            <div className="text-center group hover:scale-105 transition-all duration-300 ease-in-out">
              <div className="bg-blue-600 dark:bg-blue-500 text-white w-10 h-10 sm:w-12 sm:h-12 rounded-full flex items-center justify-center mx-auto mb-3 sm:mb-4 text-lg sm:text-xl font-bold transition-all duration-300 group-hover:scale-110 group-hover:shadow-lg">
                3
              </div>
              <h3 className="text-base sm:text-lg font-semibold text-gray-900 dark:text-white mb-2 transition-colors duration-300">
                Generate
              </h3>
              <p className="text-sm sm:text-base text-gray-600 dark:text-gray-300 transition-colors duration-300">
                Get summaries, flashcards, and quizzes
              </p>
            </div>

            <div className="text-center group hover:scale-105 transition-all duration-300 ease-in-out">
              <div className="bg-blue-600 dark:bg-blue-500 text-white w-10 h-10 sm:w-12 sm:h-12 rounded-full flex items-center justify-center mx-auto mb-3 sm:mb-4 text-lg sm:text-xl font-bold transition-all duration-300 group-hover:scale-110 group-hover:shadow-lg">
                4
              </div>
              <h3 className="text-base sm:text-lg font-semibold text-gray-900 dark:text-white mb-2 transition-colors duration-300">
                Study
              </h3>
              <p className="text-sm sm:text-base text-gray-600 dark:text-gray-300 transition-colors duration-300">
                Study smarter with AI-powered materials
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Footer */}
      <footer className="bg-gray-900 dark:bg-gray-950 text-white py-8 sm:py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <div className="flex items-center justify-center mb-3 sm:mb-4">
            <BookOpen className="h-6 w-6 sm:h-8 sm:w-8 text-blue-400" />
            <span className="ml-2 text-lg sm:text-xl font-bold">AI Study Partner</span>
          </div>
          <p className="text-sm sm:text-base text-gray-400">
            © 2025 AI Study Partner. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  );
}
