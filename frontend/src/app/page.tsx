/**
 * Homepage — Landing page with hero, features, career categories, and CTA sections.
 * Fully responsive and dark-mode compatible.
 */
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function HomePage() {
  return (
    <div className="flex flex-col">
      {/* ── Hero Section ── */}
      <section className="relative py-20 md:py-32 bg-gradient-to-br from-primary/5 via-background to-primary/10 overflow-hidden">
        {/* Decorative blurred blobs */}
        <div className="absolute -top-24 -left-24 w-96 h-96 bg-primary/10 rounded-full blur-3xl pointer-events-none" />
        <div className="absolute -bottom-24 -right-24 w-96 h-96 bg-primary/5 rounded-full blur-3xl pointer-events-none" />

        <div className="container mx-auto px-4 text-center relative z-10">
          <div className="inline-block px-4 py-1.5 mb-6 rounded-full bg-primary/10 text-primary text-sm font-medium border border-primary/20">
            🇮🇳 Built for Indian Students
          </div>
          <h1 className="text-4xl md:text-6xl font-bold tracking-tight text-foreground mb-6 leading-tight">
            Discover Your
            <span className="text-primary"> Perfect Career Path</span>
          </h1>
          <p className="text-lg md:text-xl text-muted-foreground max-w-2xl mx-auto mb-10">
            AI-powered career assessment designed specifically for Indian students.
            Get personalized recommendations based on your academics, finances, personality, and goals.
          </p>
          <div className="flex flex-col sm:flex-row gap-3 sm:gap-4 justify-center items-center">
            <Link href="/register" className="w-full sm:w-auto">
              <Button size="lg" className="text-base sm:text-lg px-6 sm:px-8 py-5 sm:py-6 w-full shadow-lg shadow-primary/25 hover:shadow-xl hover:shadow-primary/30 transition-shadow">
                Start Free Assessment →
              </Button>
            </Link>
            <Link href="/login" className="w-full sm:w-auto">
              <Button variant="outline" size="lg" className="text-base sm:text-lg px-6 sm:px-8 py-5 sm:py-6 w-full">
                Login
              </Button>
            </Link>
          </div>
        </div>
      </section>

      {/* ── Features Section ── */}
      <section className="py-12 md:py-20 bg-card">
        <div className="container mx-auto px-4">
          <h2 className="text-2xl md:text-3xl font-bold text-center mb-3">How CareerManifest Works</h2>
          <p className="text-muted-foreground text-center mb-8 md:mb-12 max-w-lg mx-auto">
            Three simple steps to discover your ideal career path
          </p>
          <div className="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto">
            <Card className="text-center p-6 hover:shadow-lg transition-shadow border-primary/10 hover:border-primary/30">
              <CardHeader>
                <div className="h-16 w-16 mx-auto mb-4 rounded-2xl bg-primary/10 flex items-center justify-center text-3xl">📝</div>
                <CardTitle className="text-lg">Answer 30 Questions</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">
                  Complete our India-focused assessment covering academics, finances, personality, and career interests.
                </p>
              </CardContent>
            </Card>

            <Card className="text-center p-6 hover:shadow-lg transition-shadow border-primary/10 hover:border-primary/30">
              <CardHeader>
                <div className="h-16 w-16 mx-auto mb-4 rounded-2xl bg-primary/10 flex items-center justify-center text-3xl">🧠</div>
                <CardTitle className="text-lg">AI Analysis</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">
                  Our weighted scoring engine evaluates your answers across 6 major career categories with India-realistic parameters.
                </p>
              </CardContent>
            </Card>

            <Card className="text-center p-6 hover:shadow-lg transition-shadow border-primary/10 hover:border-primary/30">
              <CardHeader>
                <div className="h-16 w-16 mx-auto mb-4 rounded-2xl bg-primary/10 flex items-center justify-center text-3xl">🎯</div>
                <CardTitle className="text-lg">Get Your Roadmap</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">
                  Receive your best career path, risk assessment, salary projections, required skills, and step-by-step preparation plan.
                </p>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>

      {/* ── Career Categories ── */}
      <section className="py-12 md:py-20 bg-muted/50">
        <div className="container mx-auto px-4">
          <h2 className="text-2xl md:text-3xl font-bold text-center mb-3">Career Paths We Analyze</h2>
          <p className="text-muted-foreground text-center mb-8 md:mb-12 max-w-lg mx-auto">
            Comprehensive coverage of the most popular career tracks in India
          </p>
          <div className="grid grid-cols-2 md:grid-cols-3 gap-4 md:gap-6 max-w-4xl mx-auto">
            {[
              { emoji: "💻", title: "IT / Software Jobs", desc: "TCS to Google" },
              { emoji: "📊", title: "MBA (India)", desc: "IIMs, XLRI, ISB" },
              { emoji: "🏛️", title: "Government Exams", desc: "UPSC, SSC, Banking" },
              { emoji: "🚀", title: "Startup", desc: "Build your own venture" },
              { emoji: "🎓", title: "Higher Studies (India)", desc: "GATE, M.Tech, PhD" },
              { emoji: "✈️", title: "MS Abroad", desc: "US, Canada, Europe" },
            ].map((career) => (
              <Card key={career.title} className="p-5 md:p-6 text-center hover:shadow-lg hover:border-primary/30 transition-all group">
                <div className="text-3xl md:text-4xl mb-3 group-hover:scale-110 transition-transform">{career.emoji}</div>
                <h3 className="font-semibold text-sm">{career.title}</h3>
                <p className="text-xs text-muted-foreground mt-1">{career.desc}</p>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* ── Call To Action ── */}
      <section className="py-12 md:py-20 bg-primary text-primary-foreground">
        <div className="container mx-auto px-4 text-center">
          <h2 className="text-2xl md:text-3xl font-bold mb-4 md:mb-6">Ready to Find Your Career Path?</h2>
          <p className="text-lg opacity-90 mb-8 max-w-xl mx-auto">
            Join thousands of Indian students who have discovered their ideal career through CareerManifest.
          </p>
          <Link href="/register">
            <Button size="lg" variant="secondary" className="text-lg px-8 py-6 shadow-lg hover:shadow-xl transition-shadow">
              Take the Assessment — It&apos;s Free
            </Button>
          </Link>
        </div>
      </section>
    </div>
  );
}
