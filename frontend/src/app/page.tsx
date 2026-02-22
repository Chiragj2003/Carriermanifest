/**
 * Homepage — Landing page with hero, features, career categories, and CTA sections.
 * Fully responsive with mobile-first design and dark-mode support.
 */
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function HomePage() {
  return (
    <div className="flex flex-col">
      {/* ── Hero Section ── */}
      <section className="relative py-16 sm:py-20 md:py-32 bg-gradient-to-br from-primary/5 via-background to-primary/10 overflow-hidden">
        {/* Decorative blurred blobs */}
        <div className="absolute -top-24 -left-24 w-64 sm:w-96 h-64 sm:h-96 bg-primary/10 rounded-full blur-3xl pointer-events-none" />
        <div className="absolute -bottom-24 -right-24 w-64 sm:w-96 h-64 sm:h-96 bg-primary/5 rounded-full blur-3xl pointer-events-none" />

        <div className="container mx-auto px-4 sm:px-6 text-center relative z-10">
          <div className="inline-block px-3 sm:px-4 py-1.5 mb-4 sm:mb-6 rounded-full bg-primary/10 text-primary text-xs sm:text-sm font-medium border border-primary/20">
            🇮🇳 Built for Indian Students
          </div>
          <h1 className="text-3xl sm:text-4xl md:text-6xl font-bold tracking-tight text-foreground mb-4 sm:mb-6 leading-tight">
            Discover Your
            <span className="text-primary block sm:inline"> Perfect Career Path</span>
          </h1>
          <p className="text-base sm:text-lg md:text-xl text-muted-foreground max-w-2xl mx-auto mb-8 sm:mb-10 px-2">
            AI-powered career assessment designed specifically for Indian students.
            Get personalized recommendations based on your academics, finances, personality, and goals.
          </p>
          <div className="flex flex-col sm:flex-row gap-3 sm:gap-4 justify-center items-center px-4 sm:px-0">
            <Link href="/register" className="w-full sm:w-auto">
              <Button size="lg" className="text-base sm:text-lg px-6 sm:px-8 py-5 sm:py-6 w-full shadow-lg shadow-primary/25 hover:shadow-xl hover:shadow-primary/30 transition-all active:scale-[0.98]">
                Start Free Assessment →
              </Button>
            </Link>
            <Link href="/login" className="w-full sm:w-auto">
              <Button variant="outline" size="lg" className="text-base sm:text-lg px-6 sm:px-8 py-5 sm:py-6 w-full active:scale-[0.98]">
                Login
              </Button>
            </Link>
          </div>
        </div>
      </section>

      {/* ── Stats Section (Mobile-friendly social proof) ── */}
      <section className="py-6 sm:py-8 bg-card border-b">
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-3 gap-4 sm:gap-8 max-w-3xl mx-auto text-center">
            <div>
              <p className="text-2xl sm:text-3xl font-bold text-primary">30+</p>
              <p className="text-xs sm:text-sm text-muted-foreground">Smart Questions</p>
            </div>
            <div>
              <p className="text-2xl sm:text-3xl font-bold text-primary">6</p>
              <p className="text-xs sm:text-sm text-muted-foreground">Career Paths</p>
            </div>
            <div>
              <p className="text-2xl sm:text-3xl font-bold text-primary">AI</p>
              <p className="text-xs sm:text-sm text-muted-foreground">Powered Analysis</p>
            </div>
          </div>
        </div>
      </section>

      {/* ── Features Section ── */}
      <section className="py-12 md:py-20 bg-card">
        <div className="container mx-auto px-4 sm:px-6">
          <h2 className="text-xl sm:text-2xl md:text-3xl font-bold text-center mb-2 sm:mb-3">How CareerManifest Works</h2>
          <p className="text-muted-foreground text-center text-sm sm:text-base mb-8 md:mb-12 max-w-lg mx-auto">
            Three simple steps to discover your ideal career path
          </p>
          <div className="grid sm:grid-cols-2 md:grid-cols-3 gap-4 sm:gap-6 md:gap-8 max-w-5xl mx-auto">
            <Card className="text-center p-4 sm:p-6 hover:shadow-lg transition-shadow border-primary/10 hover:border-primary/30 active:scale-[0.98] transition-transform">
              <CardHeader className="pb-2 sm:pb-4">
                <div className="h-12 w-12 sm:h-16 sm:w-16 mx-auto mb-3 sm:mb-4 rounded-2xl bg-primary/10 flex items-center justify-center text-2xl sm:text-3xl">📝</div>
                <CardTitle className="text-base sm:text-lg">Answer 30 Questions</CardTitle>
              </CardHeader>
              <CardContent className="px-2 sm:px-4">
                <p className="text-muted-foreground text-sm">
                  Complete our India-focused assessment covering academics, finances, personality, and career interests.
                </p>
              </CardContent>
            </Card>

            <Card className="text-center p-4 sm:p-6 hover:shadow-lg transition-shadow border-primary/10 hover:border-primary/30 active:scale-[0.98] transition-transform">
              <CardHeader className="pb-2 sm:pb-4">
                <div className="h-12 w-12 sm:h-16 sm:w-16 mx-auto mb-3 sm:mb-4 rounded-2xl bg-primary/10 flex items-center justify-center text-2xl sm:text-3xl">🧠</div>
                <CardTitle className="text-base sm:text-lg">AI Analysis</CardTitle>
              </CardHeader>
              <CardContent className="px-2 sm:px-4">
                <p className="text-muted-foreground text-sm">
                  Our weighted scoring engine evaluates your answers across 6 major career categories with India-realistic parameters.
                </p>
              </CardContent>
            </Card>

            <Card className="text-center p-4 sm:p-6 hover:shadow-lg transition-shadow border-primary/10 hover:border-primary/30 active:scale-[0.98] transition-transform sm:col-span-2 md:col-span-1">
              <CardHeader className="pb-2 sm:pb-4">
                <div className="h-12 w-12 sm:h-16 sm:w-16 mx-auto mb-3 sm:mb-4 rounded-2xl bg-primary/10 flex items-center justify-center text-2xl sm:text-3xl">🎯</div>
                <CardTitle className="text-base sm:text-lg">Get Your Roadmap</CardTitle>
              </CardHeader>
              <CardContent className="px-2 sm:px-4">
                <p className="text-muted-foreground text-sm">
                  Receive your best career path, risk assessment, salary projections, required skills, and step-by-step preparation plan.
                </p>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>

      {/* ── Career Categories ── */}
      <section className="py-12 md:py-20 bg-muted/50">
        <div className="container mx-auto px-4 sm:px-6">
          <h2 className="text-xl sm:text-2xl md:text-3xl font-bold text-center mb-2 sm:mb-3">Career Paths We Analyze</h2>
          <p className="text-muted-foreground text-center text-sm sm:text-base mb-8 md:mb-12 max-w-lg mx-auto">
            Comprehensive coverage of the most popular career tracks in India
          </p>
          <div className="grid grid-cols-2 md:grid-cols-3 gap-3 sm:gap-4 md:gap-6 max-w-4xl mx-auto">
            {[
              { emoji: "💻", title: "IT / Software Jobs", desc: "TCS to Google" },
              { emoji: "📊", title: "MBA (India)", desc: "IIMs, XLRI, ISB" },
              { emoji: "🏛️", title: "Government Exams", desc: "UPSC, SSC, Banking" },
              { emoji: "🚀", title: "Startup", desc: "Build your own venture" },
              { emoji: "🎓", title: "Higher Studies (India)", desc: "GATE, M.Tech, PhD" },
              { emoji: "✈️", title: "MS Abroad", desc: "US, Canada, Europe" },
            ].map((career) => (
              <Card key={career.title} className="p-4 sm:p-5 md:p-6 text-center hover:shadow-lg hover:border-primary/30 transition-all group active:scale-[0.97]">
                <div className="text-2xl sm:text-3xl md:text-4xl mb-2 sm:mb-3 group-hover:scale-110 transition-transform">{career.emoji}</div>
                <h3 className="font-semibold text-xs sm:text-sm">{career.title}</h3>
                <p className="text-[10px] sm:text-xs text-muted-foreground mt-1">{career.desc}</p>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* ── Call To Action ── */}
      <section className="py-12 md:py-20 bg-primary text-primary-foreground">
        <div className="container mx-auto px-4 sm:px-6 text-center">
          <h2 className="text-xl sm:text-2xl md:text-3xl font-bold mb-3 sm:mb-4 md:mb-6">Ready to Find Your Career Path?</h2>
          <p className="text-sm sm:text-lg opacity-90 mb-6 sm:mb-8 max-w-xl mx-auto">
            Join thousands of Indian students who have discovered their ideal career through CareerManifest.
          </p>
          <Link href="/register">
            <Button size="lg" variant="secondary" className="text-base sm:text-lg px-6 sm:px-8 py-5 sm:py-6 shadow-lg hover:shadow-xl transition-all active:scale-[0.98]">
              Take the Assessment — It&apos;s Free
            </Button>
          </Link>
        </div>
      </section>
    </div>
  );
}
