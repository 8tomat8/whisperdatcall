package main

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

const instructionPrompt = "system\nYou are an assistant in summarizing transcriptions from video calls into a list of bullet points. It must be maximum 10 bullet points."

type magic struct {
	BananaAPIKey   string `envconfig:"BANANA_API_KEY"`
	BananaModelKey string `envconfig:"BANANA_MODEL_KEY"`
	OpenAICli      *openai.Client
}

func (h magic) transcribe(ctx context.Context, filePath string) (string, error) {
	resp, err := h.OpenAICli.CreateTranscription(ctx, openai.AudioRequest{
		Model:       openai.Whisper1,
		FilePath:    filePath,
		Prompt:      "",
		Temperature: 0,
	})
	if err != nil {
		return "", errors.Wrap(err, "create transcript")
	}

	return resp.Text, nil
}

func (h magic) summarize(ctx context.Context, text string) (string, error) {
	result := make(chan string)

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	for i := 0; i < 5; i++ {
		go func() {
			resp, err := h.OpenAICli.CreateChatCompletion(
				ctx,
				openai.ChatCompletionRequest{
					Model: openai.GPT3Dot5Turbo,
					Messages: []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleSystem,
							Content: instructionPrompt,
						},
						{
							Role:    openai.ChatMessageRoleUser,
							Content: text,
						},
					},
					Stream:           false,
					Temperature:      0.3,
					TopP:             1,
					MaxTokens:        512,
					FrequencyPenalty: 0,
				},
			)

			if err != nil {
				log.Println(errors.Wrap(err, "calling openai API").Error())
			}

			if len(resp.Choices) == 0 {
				return
			}
			select {
			case <-ctx.Done():
			case result <- resp.Choices[0].Message.Content:
				cancel()
			}

			return
		}()
	}

	select {
	case text := <-result:
		return text, nil
	case <-ctx.Done():
		return "", errors.New("no result :(")
	}
	// resp, err := h.OpenAICli.CreateChatCompletion(
	// 	ctx,
	// 	openai.ChatCompletionRequest{
	// 		Model: openai.GPT3Dot5Turbo,
	// 		Messages: []openai.ChatCompletionMessage{
	// 			{
	// 				Role:    openai.ChatMessageRoleSystem,
	// 				Content: instructionPrompt,
	// 			},
	// 			{
	// 				Role:    openai.ChatMessageRoleUser,
	// 				Content: text,
	// 			},
	// 		},
	// 		Stream:           false,
	// 		Temperature:      0.7,
	// 		TopP:             1,
	// 		MaxTokens:        256,
	// 		FrequencyPenalty: 0,
	// 	},
	// )
	//
	// if err != nil {
	// 	log.Println(errors.Wrap(err, "calling openai API").Error())
	// 	return "", errors.Wrap(err, "calling openai API")
	// }
	// return resp.Choices[0].Message.Content, nil
}

const someText = `Have you ever ended a phone call with a customer and tried to replay what was said in your head? Maybe you thought there was a more effective way to describe the solution you were telling a customer about. Or maybe you forgot a key piece of information a customer shared with you about their unique case.

Call Transcription benefits and Which Call Transcription Tools you can Use on your customer service team
This is where call transcriptions come in handy. These days, call recordings and call transcriptions are an important part of an effective customer service experience. In fact, there are a number of benefits that result from call transcriptions — which we’ll review below — in addition to some tools you might choose to use on your service team or at your call center.

Let’s get started.

Get Started with HubSpot's Call Recording Software for Free
What is call transcription?
Call transcription is the process of converting conversations that take place via phone call — whether VoIP or traditional phone — into written words. (This is also known as speech-to-text transcription.) Call transcription software makes this an automatic process — and one that can happen in real-time or after a call has been recorded. Call transcriptions provide reps with scannable records of every conversation they have with customers.

Call Transcription Benefits
Here are six reasons why call transcriptions are beneficial to every business.

1. Keep records of all rep-to-customer conversations.
Similar to a call recording, a transcription is a record of your conversation that you can store for as long as you want and reference whenever you need. When you have records of rep-to-customer conversations, you can pick up with a customer wherever the last conversation ended.

You’re also able to learn about your audience, refer to key statements and highlights, share records with other members of the organization (cross-team), and analyze the success of that conversation to improve the customer experience and buyer’s journey.

2. Search and scan transcriptions for specific information.
When conversations are transcribed into written text, you’re able to efficiently scan for and reference specific highlights, keywords, and phrases. This can be helpful when telling a manger or fellow rep about a case — or if you're following up with a customer and need to reference that information. 

3. Use transcriptions for training purposes.
Transcriptions can be shared among reps and throughout the new-hire onboarding process. This is a great way to offer an example script of what a rep could potentially say (or not say) to a customer while providing support. 

4. Provide reps the opportunity to listen to their calls (and identify strengths and weaknesses).
When you stop to review your work — in this case, read call transcriptions — you’re able to identify your strengths and weaknesses. Meaning, you can pinpoint examples of what you should continue doing as well as opportunities for improvement when working with a customer.

5. Keep transcriptions for legal purposes.
Although a legal situation in which you need to provide call transcriptions and/or recordings may be rare, it's better to be prepared when it comes to the security and health of your business. By recording and transcribing calls, you have evidence and proof of what was said between reps and customers, and how situations were handled.

(Remember, all call recordings must abide by the law — be sure to review your state’s laws about call recording.)

6. Improve accessibility.
When you have and offer access to written transcripts of your calls, your conversations become accessible to everyone, including those who are hard of hearing or deaf.

Call Transcription Software
There are many call transcription software available today — some of which are strictly call transcription tools while others also assist with call recordings and other customer service or business tasks. Here are five options to get you started.

1. Fireflies.ai
fireflies.ai Call Transcription Software
Source

Fireflies.ai is an artificial intelligence (AI)-powered meeting assistant software with the ability to help record, transcribe, and search your voice conversations. The transcription feature transcribes the conversations you need it to including live meetings and past audio files.

Once you have a transcript, mark specific points in the call, leave comments, scan for keywords and highlights, and share it with team members.

Integrate Fireflies.ai with HubSpot to send meeting recordings, notes & transcripts directly to your CRM.

2. Gong
gong Call Transcription Software
Source

Gong is a revenue intelligence platform with call recording and transcription. Call transcription is automatic so you don’t have to worry about anything but the conversation at hand.

Search transcriptions for certain comments or highlights and mark the most important parts of the call. Gong analyzes call transcriptions for you and identifies key topics throughout the conversation for your records. 

Integrate Gong with HubSpot to enrich conversations with customers with account data for greater insights.

3. Jiminny
jiminny Call Transcription Software
Source

Jiminny is a revenue intelligence platform made for remote teams powered by AI — it automatically records, transcribes, and analyzes calls and meetings. After you end a call, the transcription — as well as any audio and video — is automatically sent to the cloud for easy access and analysis. 

Integrate Jiminny with HubSpot for automation and AI recording, transcription & analysis capabilities.

4. Wingman
wingman Call Transcription Software
Source

Wingman is a conversation intelligence platform that automatically transcribes calls across dialers and online/video calling and meeting tools. Wingman uses AI to analyze call recordings and transcriptions to provide insight into what your customers want and need from you. 

Integrate Wingman with HubSpot to improve rep-to-customer conversations with insights from deals.

5. Otter
otter call transcription software
Source

Otter.ai is a voice meeting notes software for all meetings, interviews, and lectures. Transcribe any conversation as well as podcasts and videos — pair any transcription or recording with meeting notes and highlights. Use Otter’s real-time transcription feature to add live captioning when speaking to a group of people at a conference.  

Call transcription has numerous benefits that have the ability to help your service team improve support and more effectively assist customers on a regular basis.`
