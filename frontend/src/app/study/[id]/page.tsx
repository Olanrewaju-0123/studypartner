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
      const flashcardsData = await apiClient.generateFlashcards(noteId);
      setFlashcards(flashcardsData);
    } catch (err) {
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
      const quizData = await apiClient.generateQuiz(noteId);
      setQuiz(quizData);
    } catch (err) {
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
            <h1 className="text-xl font-semibold text-gray-900 truncate max-w-md">
              {note.title}
            </h1>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Note Info */}
        <div className="bg-white rounded-lg shadow-sm p-6 mb-8">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-2xl font-bold text-gray-900">{note.title}</h2>
              <p className="text-gray-600 mt-1">
                {note.file_type.toUpperCase()} •{" "}
                {(note.file_size / 1024).toFixed(1)} KB • Uploaded{" "}
                {new Date(note.created_at).toLocaleDateString()}
              </p>
            </div>
          </div>
        </div>

        {/* Tabs */}
        <div className="bg-white rounded-lg shadow-sm mb-8">
          <div className="border-b border-gray-200">
            <nav className="-mb-px flex space-x-8 px-6">
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
                      flex items-center py-4 px-1 border-b-2 font-medium text-sm
                      ${
                        activeTab === tab.id
                          ? "border-blue-500 text-blue-600"
                          : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                      }
                    `}
                  >
                    <Icon className="h-5 w-5 mr-2" />
                    {tab.label}
                  </button>
                );
              })}
            </nav>
          </div>

          {/* Tab Content */}
          <div className="p-6">
            {activeTab === "summary" && (
              <div>
                {summary ? (
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">
                      Summary
                    </h3>
                    <div className="prose max-w-none">
                      <p className="text-gray-700 leading-relaxed whitespace-pre-wrap">
                        {summary.content}
                      </p>
                    </div>
                  </div>
                ) : (
                  <div className="text-center py-12">
                    <FileText className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                    <h3 className="text-lg font-medium text-gray-900 mb-2">
                      No Summary Yet
                    </h3>
                    <p className="text-gray-600 mb-6">
                      Generate an AI-powered summary of your notes
                    </p>
                    {error && (
                      <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg">
                        <p className="text-red-600 text-sm">{error}</p>
                      </div>
                    )}
                    <button
                      onClick={generateSummary}
                      disabled={generating === "summary"}
                      className="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50 flex items-center mx-auto"
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
                {flashcards.length > 0 ? (
                  <FlashcardComponent flashcards={flashcards} />
                ) : (
                  <div className="text-center py-12">
                    <Target className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                    <h3 className="text-lg font-medium text-gray-900 mb-2">
                      No Flashcards Yet
                    </h3>
                    <p className="text-gray-600 mb-6">
                      Generate interactive flashcards from your notes
                    </p>
                    <button
                      onClick={generateFlashcards}
                      disabled={generating === "flashcards"}
                      className="bg-green-600 text-white px-6 py-2 rounded-lg hover:bg-green-700 disabled:opacity-50 flex items-center mx-auto"
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
                {quiz.length > 0 ? (
                  <QuizComponent quiz={quiz} />
                ) : (
                  <div className="text-center py-12">
                    <Brain className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                    <h3 className="text-lg font-medium text-gray-900 mb-2">
                      No Quiz Yet
                    </h3>
                    <p className="text-gray-600 mb-6">
                      Generate a quiz to test your knowledge
                    </p>
                    <button
                      onClick={generateQuiz}
                      disabled={generating === "quiz"}
                      className="bg-purple-600 text-white px-6 py-2 rounded-lg hover:bg-purple-700 disabled:opacity-50 flex items-center mx-auto"
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
