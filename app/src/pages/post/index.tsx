import React, { useState } from "react"
import { useForm } from "react-hook-form";
import {
	FormControl,
	FormLabel,
	Input,
	VStack,
	Container,
	Stack,
	Button,
} from '@chakra-ui/react';
import Template from "../template";
import {sessionInformation} from '../../../env'
type quizType = {
	title: string
	answer: string
	wrongAnswer1: string
	wrongAnswer2: string
	wrongAnswer3: string
	explanation: string
}
const quizData: quizType = {
	title: "",
	answer: "",
	wrongAnswer1: "",
	wrongAnswer2: "",
	wrongAnswer3: "",
	explanation: "",
}
const Posts = () => {
	const { register, handleSubmit, formState, formState: { errors }, getValues } = useForm({
		mode: "onTouched",
	});
	const [posts, setPosts] = useState([])
	const [randomAnswer, setRandomAnswer] = useState([])
	const [answer, setAnswer] = useState<string>()
	const setData = () => {
		const hasData = getValues(["title", "answer", "wrongAnswer1", "wrongAnswer2", "wrongAnswer3", "explanation"])
		quizData.title = hasData[0]
		quizData.answer = hasData[1]
		quizData.wrongAnswer1 = hasData[2]
		quizData.wrongAnswer2 = hasData[3]
		quizData.wrongAnswer3 = hasData[4]
		quizData.explanation = hasData[5]
	}
	const RandomAnswer = () => {
		let answer = [quizData.answer, quizData.wrongAnswer1, quizData.wrongAnswer2, quizData.wrongAnswer3]
		for (let i = answer.length - 1; i > 0; i--) {
			let j = Math.floor(Math.random() * (i + 1));
			let tmp = answer[i];
			answer[i] = answer[j];
			answer[j] = tmp;
		}
		setRandomAnswer(answer)
	}
	const test = () => {
		setData()
		RandomAnswer()
	}
	const getAnswer = (event) => {
		if (quizData.answer == event.target.value) {
			setAnswer("正解！")
		} else {
			setAnswer("不正解！")
		}
	}
	const ApiFetch = () => {
		setData()
		fetch(`${sessionInformation.backendHost}/post`, {
			mode: "cors",
			method: "POST",
			headers: { "Content-Type": "application/json", }, // JSON形式のデータのヘッダー
			credentials: 'include',
			body: JSON.stringify(quizData),
		})
			.then((res) => res.json())
			.then((data) => {
				setPosts(data);
			})
			.catch((err) => { console.log(err) })
	};

	return (<>
		<Template />
		<FormControl onSubmit={handleSubmit(setData)}
			isInvalid={errors.title ? true : false}>
			<FormLabel>問題文</FormLabel>
			<Input
				type="string"
				placeholder="例:この中でフロントエンド言語はどれ？"
				{...register("title", {
					required: true,
					pattern: {
						value: /^[^^＾"”`‘'’<>＜＞_＿%$#＆％＄|￥]{1,200}$/,
						message: '特殊文字を使用しないでください' // JS only: <p>error message</p> TS only support string
					}
				})}
			/>
			{errors.title && errors.title.message}
		</FormControl>
		<FormControl onSubmit={handleSubmit(setData)}
			isInvalid={errors.answer ? true : false}>
			<FormLabel>正答</FormLabel>
			<Input
				type="body"
				placeholder="例:JavaScript"
				{...register("answer", {
					required: true,
					pattern: {
						value: /^[^^＾"”`‘'’<>＜＞_＿%$#＆％＄|￥]{1,200}$/,
						message: '特殊文字を使用しないでください' // JS only: <p>error message</p> TS only support string
					}
				})}
			/>
			{errors.answer && errors.answer.message}
		</FormControl>

		<FormControl onSubmit={handleSubmit(setData)}
			isInvalid={errors.wrongAnswer1 ? true : false}>
			<FormLabel>誤答</FormLabel>
			<Input
				type="body"
				placeholder="例:Go"
				{...register("wrongAnswer1", {
					required: true,
					pattern: {
						value: /^[^^＾"”`‘'’<>＜＞_＿%$#＆％＄|￥]{1,200}$/,
						message: '特殊文字を使用しないでください' // JS only: <p>error message</p> TS only support string
					}
				})}
			/>
			{errors.wrongAnswer1 && errors.wrongAnswer1.message}
		</FormControl>

		<FormControl onSubmit={handleSubmit(setData)}
			isInvalid={errors.wrongAnswer2 ? true : false}>
			<FormLabel>誤答</FormLabel>
			<Input
				type="body"
				placeholder="例:PHP"
				{...register("wrongAnswer2", {
					required: true,
					pattern: {
						value: /^[^^＾"”`‘'’<>＜＞_＿%$#＆％＄|￥]{1,200}$/,
						message: '特殊文字を使用しないでください' // JS only: <p>error message</p> TS only support string
					}
				})}
			/>
			{errors.wrongAnswer2 && errors.wrongAnswer2.message}
		</FormControl>

		<FormControl onSubmit={handleSubmit(setData)}
			isInvalid={errors.wrongAnswer3 ? true : false}>
			<FormLabel>誤答</FormLabel>
			<Input
				type="body"
				placeholder="例:Ruby"
				{...register("wrongAnswer3", {
					required: true,
					pattern: {
						value: /^[^^＾"”`‘'’<>＜＞_＿%$#＆％＄|￥]{1,200}$/,
						message: '特殊文字を使用しないでください' // JS only: <p>error message</p> TS only support string
					}
				})}
			/>
			{errors.wrongAnswer3 && errors.wrongAnswer3.message}
		</FormControl>
		<Stack>
			<FormControl onSubmit={handleSubmit(setData)}
				isInvalid={errors.explanation ? true : false}>
				<FormLabel>解説文</FormLabel>
				<Input
					type="body"
					placeholder="JavaScriptだけがフロントエンド言語だよ！"
					{...register("explanation", {
						required: true,
						pattern: {
							value: /^[^^＾"”`‘'’<>＜＞_＿%$#＆％＄|￥]{1,200}$/,
							message: '特殊文字を使用しないでください' // JS only: <p>error message</p> TS only support string
						}
					})}
				/>
				{errors.explanation && errors.explanation.message}
			</FormControl>
		</Stack>
		<Button type="submit"
			colorScheme="teal"
			onClick={ApiFetch}
			disabled={!formState.isValid}
		>送信</Button>
		<Stack>
			<Button onClick={test}>プレビュー</Button>
		</Stack>
		<Stack>
			<>問題文:{quizData.title}</>
		</Stack>
		<VStack>
			<Container>回答1:{randomAnswer[0]}<Button margin="2" onClick={getAnswer} value={randomAnswer[0]}>これにする</Button></Container>
			<Container>回答2:{randomAnswer[1]}<Button margin="2" onClick={getAnswer} value={randomAnswer[1]}>これにする</Button></Container>
			<Container>回答3:{randomAnswer[2]}<Button margin="2" onClick={getAnswer} value={randomAnswer[2]}>これにする</Button></Container>
			<Container>回答4:{randomAnswer[3]}<Button margin="2" onClick={getAnswer} value={randomAnswer[3]}>これにする</Button></Container>
		</VStack>
		<Stack>
			<>{answer}</>
			{answer && <>正解は{quizData.answer}です</>}
		</Stack>
		<Stack>
			{answer && <>解説文{quizData.explanation}</>}
		</Stack>
		{JSON.stringify(posts)}
	</>)
}
export default Posts