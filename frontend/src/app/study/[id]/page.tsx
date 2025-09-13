"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import {
  ArrowLeft,
  FileText,
  Brain,
  Target,
  Search,
  Loader2,
} from "lucide-react";
import Link from "next/link";
import { apiClient } from "@/utils/api";
import { Note, Summary, Flashcard, Quiz } from "@/types";
import FlashcardComponent from "@/components/Flashcard";
import QuizComponent from "@/components/Quiz";
import ThemeToggle from "@/components/ThemeToggle";

export default function StudyPage() {
  const params = useParams();
  const router = useRouter();
  const noteId = parseInt(params.id as string);

  const [note, setNote] = useState<Note | null>(null);
  const [summary, setSummary] = useState<Summary | null>(null);
  const [flashcards, setFlashcards] = useState<Flashcard[]>([]);
  const [quiz, setQuiz] = useState<Quiz[]>([]);
  const [activeTab, setActiveTab] = useState<"summary" | "flashcards" | "quiz">(
    "summary"
  );
  const [loading, setLoading] = useState(true);
  const [generating, setGenerating] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadNote();
  }, [noteId]);

  const loadNote = async () => {
    try {
      setLoading(true);
      const noteData = await apiClient.getNote(noteId);
      setNote(noteData);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load note");
    } finally {
      setLoading(false);
    }
  };

  const generateSummary = async () => {
    try {
      setGenerating("summary");
      setError(null); // Clear any previous errors
      const summaryData = await apiClient.generateSummary(noteId);
      setSummary(summaryData);
    } catch (err) {
      console.error("Summary generation error:", err);
      setError(
        err instanceof Error ? err.message : "Failed to generate summary"
      );
    } finally {
      setGenerating(null);
    }
  };

  const generateFlashcards = async () => {
    try {
      setGenerating("flashcards");
      setError(null); // Clear any previous errors
      const flashcardsData = await apiClient.generateFlashcards(noteId);
      console.log("Generated flashcards:", flashcardsData);
      setFlashcards(flashcardsData || []);
    } catch (err) {
      console.error("Flashcard generation error:", err);
      setError(
        err instanceof Error ? err.message : "Failed to generate flashcards"
      );
    } finally {
      setGenerating(null);
    }
  };

  const generateQuiz = async () => {
    try {
      setGenerating("quiz");
      setError(null); // Clear any previous errors
      const quizData = await apiClient.generateQuiz(noteId);
      console.log("Generated quiz:", quizData);
      setQuiz(quizData || []);
    } catch (err) {
      console.error("Quiz generation error:", err);
      setError(err instanceof Error ? err.message : "Failed to generate quiz");
    } finally {
      setGenerating(null);
    }
  };

  const loadExistingContent = async (type: string) => {
    try {
      if (type === "summary") {
        const summaryData = await apiClient.getSummary(noteId);
        setSummary(summaryData);
      } else if (type === "flashcards") {
        const flashcardsData = await apiClient.getFlashcards(noteId);
        setFlashcards(flashcardsData);
      } else if (type === "quiz") {
        const quizData = await apiClient.getQuiz(noteId);
        setQuiz(quizData);
      }
    } catch (err) {
      // Content doesn't exist yet, will need to generate
    }
  };

  const handleTabChange = async (tab: "summary" | "flashcards" | "quiz") => {
    setActiveTab(tab);
    await loadExistingContent(tab);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <Loader2 className="h-12 w-12 text-blue-500 animate-spin mx-auto mb-4" />
          <p className="text-gray-600">Loading note...</p>
        </div>
      </div>
    );
  }

  if (error || !note) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-600 mb-4">{error || "Note not found"}</p>
          <Link href="/" className="text-blue-600 hover:text-blue-700">
            Return to Home
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      {/* Header */}
      <div className="bg-white dark:bg-gray-800 shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center min-w-0 flex-1">
              <Link
                href="/"
                className="flex items-center text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white text-sm sm:text-base"
              >
                <ArrowLeft className="h-4 w-4 sm:h-5 sm:w-5 mr-1 sm:mr-2" />
                <span className="hidden sm:inline">Back to Home</span>
                <span className="sm:hidden">Back</span>
              </Link>
            </div>
            <h1 className="text-lg sm:text-xl font-semibold text-gray-900 dark:text-white truncate mx-2 sm:mx-4 max-w-xs sm:max-w-md">
              {note.title}
            </h1>
            <div className="flex items-center">
              <ThemeToggle />
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
        {/* Note Info */}
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-4 sm:p-6 mb-6 sm:mb-8 hover:shadow-md transition-all duration-300 ease-in-out">
          <div className="flex items-center justify-between">
            <div className="min-w-0 flex-1">
              <h2 className="text-xl sm:text-2xl font-bold text-gray-900 dark:text-white truncate transition-colors duration-300">{note.title}</h2>
              <p className="text-sm sm:text-base text-gray-600 dark:text-gray-300 mt-1 transition-colors duration-300">
                {note.file_type.toUpperCase()} •{" "}
                {(note.file_size / 1024).toFixed(1)} KB • Uploaded{" "}
                {new Date(note.created_at).toLocaleDateString()}
              </p>
            </div>
          </div>
        </div>

        {/* Tabs */}
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm mb-6 sm:mb-8 hover:shadow-md transition-all duration-300 ease-in-out">
          <div className="border-b border-gray-200 dark:border-gray-700">
            <nav className="-mb-px flex space-x-4 sm:space-x-8 px-4 sm:px-6 overflow-x-auto">
              {[
                { id: "summary", label: "Summary", icon: FileText },
                { id: "flashcards", label: "Flashcards", icon: Target },
                { id: "quiz", label: "Quiz", icon: Brain },
              ].map((tab) => {
                const Icon = tab.icon;
                return (
                  <button
                    key={tab.id}
                    onClick={() => handleTabChange(tab.id as any)}
                    className={`
                      flex items-center py-3 sm:py-4 px-1 border-b-2 font-medium text-xs sm:text-sm whitespace-nowrap transition-all duration-300 ease-in-out
                      ${
                        activeTab === tab.id
                          ? "border-blue-500 dark:border-blue-400 text-blue-600 dark:text-blue-400"
                          : "border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600"
                      }
                    `}
                  >
                    <Icon className="h-4 w-4 sm:h-5 sm:w-5 mr-1 sm:mr-2" />
                    {tab.label}
                  </button>
                );
              })}
            </nav>
          </div>

          {/* Tab Content */}
          <div className="p-4 sm:p-6">
            {activeTab === "summary" && (
              <div>
                {summary ? (
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
                      Summary
                    </h3>
                    <div className="prose max-w-none">
                      <p className="text-gray-700 dark:text-gray-300 leading-relaxed whitespace-pre-wrap">
                        {summary.content}
                      </p>
                    </div>
                  </div>
                ) : (
                  <div className="text-center py-12">
                    <FileText className="h-12 w-12 text-gray-400 dark:text-gray-500 mx-auto mb-4" />
                    <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">
                      No Summary Yet
                    </h3>
                    <p className="text-gray-600 dark:text-gray-300 mb-6">
                      Generate an AI-powered summary of your notes
                    </p>
                    {error && (
                      <div className="mb-4 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
                        <p className="text-red-600 dark:text-red-400 text-sm">{error}</p>
                      </div>
                    )}
                    <button
                      onClick={generateSummary}
                      disabled={generating === "summary"}
                      className="bg-blue-600 dark:bg-blue-500 text-white px-6 py-2 rounded-lg hover:bg-blue-700 dark:hover:bg-blue-600 disabled:opacity-50 flex items-center mx-auto"
                    >
                      {generating === "summary" ? (
                        <>
                          <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                          Generating...
                        </>
                      ) : (
                        <>
                          <Brain className="h-4 w-4 mr-2" />
                          Generate Summary
                        </>
                      )}
                    </button>
                  </div>
                )}
              </div>
            )}

            {activeTab === "flashcards" && (
              <div>
                {flashcards && flashcards.length > 0 ? (
                  <FlashcardComponent flashcards={flashcards} />
                ) : (
                  <div className="text-center py-12">
                    <Target className="h-12 w-12 text-gray-400 dark:text-gray-500 mx-auto mb-4" />
                    <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">
                      No Flashcards Yet
                    </h3>
                    <p className="text-gray-600 dark:text-gray-300 mb-6">
                      Generate interactive flashcards from your notes
                    </p>
                    {error && (
                      <div className="mb-4 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
                        <p className="text-red-600 dark:text-red-400 text-sm">{error}</p>
                      </div>
                    )}
                    <button
                      onClick={generateFlashcards}
                      disabled={generating === "flashcards"}
                      className="bg-green-600 dark:bg-green-500 text-white px-6 py-2 rounded-lg hover:bg-green-700 dark:hover:bg-green-600 disabled:opacity-50 flex items-center mx-auto"
                    >
                      {generating === "flashcards" ? (
                        <>
                          <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                          Generating...
                        </>
                      ) : (
                        <>
                          <Target className="h-4 w-4 mr-2" />
                          Generate Flashcards
                        </>
                      )}
                    </button>
                  </div>
                )}
              </div>
            )}

            {activeTab === "quiz" && (
              <div>
                {quiz && quiz.length > 0 ? (
                  <QuizComponent quiz={quiz} />
                ) : (
                  <div className="text-center py-12">
                    <Brain className="h-12 w-12 text-gray-400 dark:text-gray-500 mx-auto mb-4" />
                    <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">
                      No Quiz Yet
                    </h3>
                    <p className="text-gray-600 dark:text-gray-300 mb-6">
                      Generate a quiz to test your knowledge
                    </p>
                    {error && (
                      <div className="mb-4 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
                        <p className="text-red-600 dark:text-red-400 text-sm">{error}</p>
                      </div>
                    )}
                    <button
                      onClick={generateQuiz}
                      disabled={generating === "quiz"}
                      className="bg-purple-600 dark:bg-purple-500 text-white px-6 py-2 rounded-lg hover:bg-purple-700 dark:hover:bg-purple-600 disabled:opacity-50 flex items-center mx-auto"
                    >
                      {generating === "quiz" ? (
                        <>
                          <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                          Generating...
                        </>
                      ) : (
                        <>
                          <Brain className="h-4 w-4 mr-2" />
                          Generate Quiz
                        </>
                      )}
                    </button>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
