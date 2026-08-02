import { getTranslations } from "next-intl/server";
import { QuestionsMeClient } from "./QuestionsMeClient";

export async function generateMetadata() {
  const t = await getTranslations("Questions");
  return { title: t("my_questions") };
}

export default function QuestionsMePage() {
  return <QuestionsMeClient />;
}
